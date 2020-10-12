package config

import (
	"fmt"
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
			p, err := ConvertMapKey(v)
			if err != nil {
				return fmt.Errorf("convert entry.params: entry_name: %s: key: %s: %w", entry.Name, k, err)
			}
			entry.Params[k] = p
		}
		for j, test := range entry.Tests {
			p, err := ConvertMapKey(test.Equal)
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

func ConvertMapKey(data interface{}) (interface{}, error) {
	switch t := data.(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(t))
		for k, v := range t {
			s, ok := k.(string)
			if !ok {
				return nil, fmt.Errorf("the map key should be string: %+v", k)
			}
			val, err := ConvertMapKey(v)
			if err != nil {
				return nil, fmt.Errorf("key: %s: %w", s, err)
			}
			m[s] = val
		}
		return m, nil
	case []interface{}:
		for i, v := range t {
			val, err := ConvertMapKey(v)
			if err != nil {
				return nil, fmt.Errorf("index: %d: %w", i, err)
			}
			t[i] = val
		}
		return t, nil
	default:
		return data, nil
	}
}
