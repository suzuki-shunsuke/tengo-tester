package cli

import (
	"context"
	"io"

	"github.com/suzuki-shunsuke/tengo-tester/pkg/constant"
	"github.com/urfave/cli/v2"
)

type Runner struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

func (runner Runner) Run(ctx context.Context, args ...string) error {
	app := cli.App{
		Name:    "tengo-tester",
		Usage:   "test tengo scripts. https://github.com/suzuki-shunsuke/tengo-tester",
		Version: constant.Version,
		Commands: []*cli.Command{
			{
				Name:   "run",
				Usage:  "test tengo scripts",
				Action: runner.action,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "log-level",
						Usage: "log level",
					},
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Usage:   "configuration file path",
					},
				},
			},
			{
				Name:   "init",
				Usage:  "generate a configuration file if it doesn't exist",
				Action: runner.initAction,
			},
		},
	}

	return app.RunContext(ctx, args)
}
