package config

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
