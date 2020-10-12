package cli

import (
	"io/ioutil"
	"os"

	"github.com/urfave/cli/v2"
)

const cfgTpl = `---
# Configuration file of tengo-tester, which is a CLI tool to test tengo scripts
# https://github.com/suzuki-shunsuke/tengo-tester
entries:
- name: main
  script_file: foo.tengo
  params: {}
  tests:
  - name: test
    script: "" # tengo script to evaluate the result
  - name: test 2
    script_file: foo_test.tengo
  - name: test 3
    var_name: foo
    equal: "foo"
`

func (runner Runner) initAction(c *cli.Context) error {
	if _, err := os.Stat(".tengo-tester.yml"); err == nil {
		return nil
	}
	if _, err := os.Stat(".tengo-tester.yaml"); err == nil {
		return nil
	}
	if err := ioutil.WriteFile(".tengo-tester.yaml", []byte(cfgTpl), 0o755); err != nil { //nolint:gosec
		return err
	}
	return nil
}
