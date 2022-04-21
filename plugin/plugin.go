package plugin

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/ergoapi/util/file"
	"github.com/sirupsen/logrus"
)

type (
	Ext struct {
		Debug bool
	}
	Push struct {
		Username string
		Password string
		Token    string
		Hub      string
		Context  string // charts directory
		Multi    bool   // multi-charts upload
		Force    bool   // force upload
	}
	Plugin struct {
		Ext  Ext
		Push Push
	}
)

type cmd struct {
	pack *exec.Cmd
	push *exec.Cmd
	name string
}

func (p Plugin) Exec() error {
	var cmds []*cmd
	p.Push.Context = strings.TrimSuffix(p.Push.Context, "/")
	if p.Push.Multi {
		logrus.Debugf("multi-charts upload: %s\n", p.Push.Context)
		if file.CheckFileExists(fmt.Sprintf("%s/%s", p.Push.Context, "Chart.yaml")) {
			logrus.Warnf("found %s/Chart.yaml, not multi-charts will only upload current context charts\n", p.Push.Context)
			cmds = append(cmds, p.pushAction(p.Push.Context))
		} else {
			charts, err := file.DirFilesList(p.Push.Context, "Chart.yaml", "")
			if err != nil {
				return err
			}
			for _, chart := range charts {
				cmds = append(cmds, p.pushAction(chart))
			}
		}
	} else {
		cmds = append(cmds, p.pushAction(p.Push.Context))
	}
	for _, cmd := range cmds {
		var b bytes.Buffer
		cmd.push.Stdout = os.Stdout
		cmd.push.Stderr = &b
		if p.Ext.Debug {
			trace(cmd.push)
		}
		err := cmd.push.Run()
		if err != nil {
			if error409(b.String()) {
				logrus.Warnf("upload [%s] chart already exists, skip\n", cmd.name)
			} else {
				logrus.Errorf("upload [%s] chart failed\n", cmd.name)
			}
		} else {
			logrus.Infof("upload [%s] chart success\n", cmd.name)
		}
	}
	return nil
}

func (p Plugin) pushAction(path string) *cmd {
	force := ""
	if p.Push.Force {
		force = "--force"
	}
	var chartname string
	var chartpath string
	// #nosec G204
	t := strings.Split(path, "/")
	if strings.HasSuffix(path, "Chart.yaml") {
		chartname = t[len(t)-2]
		chartpath = strings.TrimSuffix(fmt.Sprintf("%s/%s", p.Push.Context, path), "/Chart.yaml")
	} else {
		chartname = strings.Trim(t[len(t)-1], "/")
		chartpath = path
	}
	if len(p.Push.Token) == 0 && (len(p.Push.Username) == 0 || len(p.Push.Password) == 0) {
		return &cmd{
			// #nosec
			push: exec.Command("helm", "cm-push", chartpath, p.Push.Hub, force),
			name: chartname,
		}
	} else if len(p.Push.Token) > 0 {
		return &cmd{
			// #nosec
			push: exec.Command("helm", "cm-push", chartpath, p.Push.Hub, "--access-token", p.Push.Token, force),
			name: chartname,
		}
	}
	return &cmd{
		// #nosec
		push: exec.Command("helm", "cm-push", chartpath, p.Push.Hub, "--username", p.Push.Username, "--password", p.Push.Password, force),
		// pack: exec.Command("helm", "package", chartpath),
		name: chartname,
	}
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func trace(cmd *exec.Cmd) {
	fmt.Fprintf(os.Stdout, "+ %s\n", strings.Join(cmd.Args, " "))
}

func error409(err string) bool {
	if strings.Contains(err, "409") || strings.Contains(err, "exists") {
		return true
	}
	return false
}
