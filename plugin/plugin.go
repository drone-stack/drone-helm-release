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
		Context  string   // charts directory
		Multi    bool     // multi-charts upload
		Force    bool     // force upload
		Exthub   []string // common hub
	}
	Plugin struct {
		Ext  Ext
		Push Push
	}
)

type cmd struct {
	build    *exec.Cmd
	pack     *exec.Cmd
	push     *exec.Cmd
	name     string
	depchart string
	path     string
}

func (p Plugin) Exec() error {
	p.prepareRepoAdd()
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
		if err := p.build(cmd); err != nil {
			continue
		}
		if err := p.pack(cmd); err != nil {
			continue
		}
		cmd.push.Stdout = os.Stdout
		cmd.push.Stderr = &b
		cmd.push.Dir = cmd.path
		if p.Ext.Debug {
			p.trace(cmd.push)
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
		os.RemoveAll(cmd.depchart)
	}
	return nil
}

func (p Plugin) build(c *cmd) error {
	c.build.Dir = c.path
	if p.Ext.Debug {
		c.build.Stdout = os.Stdout
		c.build.Stderr = os.Stderr
		p.trace(c.build)
	}
	if err := c.build.Run(); err != nil {
		logrus.Warnf("helm build [%s] failed: %v\n", c.name, err)
		return err
	}
	logrus.Debugf("helm build [%s] success\n", c.name)
	return nil
}

func (p Plugin) pack(c *cmd) error {
	c.pack.Dir = c.path
	if p.Ext.Debug {
		c.pack.Stdout = os.Stdout
		c.pack.Stderr = os.Stderr
		p.trace(c.pack)
	}
	if err := c.pack.Run(); err != nil {
		logrus.Warnf("helm package [%s] failed: %v\n", c.name, err)
		return err
	}
	logrus.Debugf("helm package [%s] success\n", c.name)
	return nil
}

func (p Plugin) pushAction(path string) *cmd {
	force := ""
	if p.Push.Force {
		force = "--force"
	}
	var chartpath string
	cmdmeta := cmd{}
	// #nosec G204
	t := strings.Split(path, "/")
	if strings.HasSuffix(path, "Chart.yaml") {
		cmdmeta.name = t[len(t)-2]
		chartpath = strings.TrimSuffix(fmt.Sprintf("%s/%s", p.Push.Context, path), "/Chart.yaml")
	} else {
		cmdmeta.name = strings.Trim(t[len(t)-1], "/")
		chartpath = path
	}
	cmdmeta.path = chartpath
	cmdmeta.depchart = fmt.Sprintf("%s/charts", chartpath)
	// #nosec
	cmdmeta.build = exec.Command("helm", "dependency", "build")
	// #nosec
	cmdmeta.pack = exec.Command("helm", "package", ".", "-u")
	//cmdmeta.pack = exec.Command("helm", "package", "-u", chartpath, "-d", cmdmeta.depchart)
	if len(p.Push.Token) == 0 && (len(p.Push.Username) == 0 || len(p.Push.Password) == 0) {
		// #nosec
		cmdmeta.push = exec.Command("helm", "cm-push", ".", p.Push.Hub, force)
	} else if len(p.Push.Token) > 0 {
		// #nosec
		cmdmeta.push = exec.Command("helm", "cm-push", ".", p.Push.Hub, "--access-token", p.Push.Token, force)
	} else {
		// #nosec
		cmdmeta.push = exec.Command("helm", "cm-push", ".", p.Push.Hub, "--username", p.Push.Username, "--password", p.Push.Password, force)
	}
	return &cmdmeta
}

func (p Plugin) prepareRepoAdd() {
	// #nosec
	if err := p.runCmd("repo", "add", "hc-default", p.Push.Hub); err != nil {
		logrus.Errorf("helm repo add hc-default failed: %v", err)
	}
	for i, hub := range p.Push.Exthub {
		// #nosec
		if err := p.runCmd("repo", "add", fmt.Sprintf("hc-%d", i), hub); err != nil {
			logrus.Errorf("helm repo add hc-%d failed: %v", i, err)
		}
	}
}

func (p Plugin) runCmd(args ...string) error {
	// #nosec
	cmd := exec.Command("helm", args...)
	cmd.Stdout = os.Stdout
	if p.Ext.Debug {
		cmd.Stderr = os.Stderr
		p.trace(cmd)
	}
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

// trace writes each command to stdout with the command wrapped in an xml
// tag so that it can be extracted and displayed in the logs.
func (p Plugin) trace(cmd *exec.Cmd) {
	key := strings.Join(cmd.Args, " ")
	if len(p.Push.Token) > 0 {
		key = strings.ReplaceAll(key, p.Push.Token, "******")
	}

	if len(p.Push.Username) > 0 {
		key = strings.ReplaceAll(key, p.Push.Username, "******")
	}

	if len(p.Push.Password) > 0 {
		key = strings.ReplaceAll(key, p.Push.Password, "******")
	}

	fmt.Fprintf(os.Stdout, "+ %s\n", key)
}

func error409(err string) bool {
	if strings.Contains(err, "409") || strings.Contains(err, "exists") {
		return true
	}
	return false
}
