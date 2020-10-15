package config

import (
	"fmt"

	"github.com/suzuki-shunsuke/go-convmap/convmap"
)

type Config struct {
	Entries  []Entry
	LogLevel string `yaml:"log_level"`
}

type Entry struct {
	Name       string
	ScriptFile string `yaml:"script_file"`
	Params     map[string]interface{}
	Tests      []Test
}

type Test struct {
	Name       string
	Script     string
	ScriptFile string `yaml:"script_file"`
	VarName    string `yaml:"var_name"`
	Equal      interface{}
	IsNil      bool `yaml:"is_nil"`
	IsNotNil   bool `yaml:"is_not_nil"`
}

func ConvertConfig(cfg Config) error {
	for i, entry := range cfg.Entries {
		for k, v := range entry.Params {
			p, err := convmap.Convert(v)
			if err != nil {
				return fmt.Errorf("convert entry.params: entry_name: %s: key: %s: %w", entry.Name, k, err)
			}
			entry.Params[k] = p
		}
		for j, test := range entry.Tests {
			p, err := convmap.Convert(test.Equal)
			if err != nil {
				return fmt.Errorf("convert test.Equal: entry_name: %s: test_name: %s: %w", entry.Name, test.Name, err)
			}
			test.Equal = p
			entry.Tests[j] = test
		}
		cfg.Entries[i] = entry
	}
	return nil
}
