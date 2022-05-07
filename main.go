package main

import (
	"os"

	"github.com/drone-stack/drone-plugin-template/plugin"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var (
	version = "0.0.5"
)

type formatter struct{}

func (*formatter) Format(entry *logrus.Entry) ([]byte, error) {
	return []byte(entry.Message), nil
}

func init() {
	logrus.SetFormatter(new(formatter))
	logrus.SetOutput(os.Stdout)
	logrus.SetLevel(logrus.InfoLevel)
}

func main() {
	// Load env-file if it exists first
	if env := os.Getenv("PLUGIN_ENV_FILE"); env != "" {
		_ = godotenv.Load(env)
	}

	app := cli.NewApp()
	app.Name = "drone helm release plugin"
	app.Usage = "drone helm release plugin"
	app.Action = run
	app.Version = version
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			Usage:  "debug",
			EnvVar: "PLUGIN_DEBUG",
		},
		cli.StringFlag{
			Name:   "username",
			Usage:  "username",
			EnvVar: "PLUGIN_USERNAME",
		},
		cli.StringFlag{
			Name:   "password",
			Usage:  "password",
			EnvVar: "PLUGIN_PASSWORD",
		},
		cli.StringFlag{
			Name:   "token",
			Usage:  "token",
			EnvVar: "PLUGIN_TOKEN",
		},
		cli.StringFlag{
			Name:     "hub",
			Usage:    "hub",
			EnvVar:   "PLUGIN_HUB",
			Required: true,
		},
		cli.StringFlag{
			Name:   "context",
			Usage:  "context",
			EnvVar: "PLUGIN_CONTEXT",
			Value:  ".",
		},
		cli.BoolFlag{
			Name:   "multi",
			Usage:  "multi",
			EnvVar: "PLUGIN_MULTI",
		},
		cli.BoolFlag{
			Name:   "force",
			Usage:  "force",
			EnvVar: "PLUGIN_FORCE",
		},
		cli.StringSliceFlag{
			Name:   "exthub",
			Usage:  "exthub",
			EnvVar: "PLUGIN_EXTHUB",
		},
		cli.StringFlag{
			Name:   "exclude",
			Usage:  "exclude",
			EnvVar: "PLUGIN_EXCLUDE",
		},
	}

	if err := app.Run(os.Args); err != nil {
		logrus.Fatal(err)
	}
}

func run(c *cli.Context) error {
	plugin := plugin.Plugin{
		Ext: plugin.Ext{
			Debug: c.Bool("debug"),
		},
		Push: plugin.Push{
			Username: c.String("username"),
			Password: c.String("password"),
			Token:    c.String("token"),
			Hub:      c.String("hub"),
			Context:  c.String("context"),
			Multi:    c.Bool("multi"),
			Force:    c.Bool("force"),
			Exthub:   c.StringSlice("exthub"),
			Exclude:  c.String("exclude"),
		},
	}

	if plugin.Ext.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err := plugin.Exec(); err != nil {
		logrus.Fatal(err)
	}
	return nil
}
