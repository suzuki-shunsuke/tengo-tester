package cli

import (
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/suzuki-shunsuke/go-findconfig/findconfig"
	"github.com/suzuki-shunsuke/tengo-tester/pkg/config"
	"github.com/suzuki-shunsuke/tengo-tester/pkg/controller"
	"github.com/suzuki-shunsuke/tengo-tester/pkg/file"
	"github.com/urfave/cli/v2"
)

func (runner Runner) setCLIArg(c *cli.Context, cfg config.Config) config.Config {
	if logLevel := c.String("log-level"); logLevel != "" {
		cfg.LogLevel = logLevel
	}
	return cfg
}

func (runner Runner) action(c *cli.Context) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	reader := config.Reader{
		ExistFile: findconfig.Exist,
	}
	cfg, cfgPath, err := reader.FindAndRead(c.String("config"), wd)
	if err != nil {
		return err
	}

	cfg = runner.setCLIArg(c, cfg)

	if cfg.LogLevel != "" {
		lvl, err := logrus.ParseLevel(cfg.LogLevel)
		if err != nil {
			logrus.WithFields(logrus.Fields{
				"log_level": cfg.LogLevel,
			}).WithError(err).Error("the log level is invalid")
		}
		logrus.SetLevel(lvl)
	}

	logrus.WithFields(logrus.Fields{
		"log_level": cfg.LogLevel,
	}).Debug("config")

	ctrl := controller.Controller{
		Config:     cfg,
		FileReader: file.Reader{},
	}

	return ctrl.Run(c.Context, filepath.Dir(cfgPath))
}
