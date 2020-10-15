package expr

import (
	"context"

	"github.com/d5/tengo/v2"
	"github.com/d5/tengo/v2/stdlib"
)

type Program struct {
	source string
}

func New(expression string) Program {
	return Program{
		source: expression,
	}
}

func (prog Program) Run(params map[string]interface{}) (map[string]interface{}, error) {
	if prog.source == "" {
		return nil, nil
	}
	script := tengo.NewScript([]byte(prog.source))
	script.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
	for k, v := range params {
		if err := script.Add(k, v); err != nil {
			return nil, err
		}
	}
	compiled, err := script.RunContext(context.Background())
	if err != nil {
		return nil, err
	}
	vars := compiled.GetAll()
	m := make(map[string]interface{}, len(vars))
	for _, v := range vars {
		m[v.Name()] = convertInt64ToInt(v.Value())
	}

	return m, nil
}

func convertInt64ToInt(data interface{}) interface{} {
	switch t := data.(type) {
	case map[string]interface{}:
		m := make(map[string]interface{}, len(t))
		for k, v := range t {
			m[k] = convertInt64ToInt(v)
		}
		return m
	case []interface{}:
		arr := make([]interface{}, len(t))
		for i, v := range t {
			arr[i] = convertInt64ToInt(v)
		}
		return arr
	case int64:
		return int(t)
	default:
		return data
	}
}
