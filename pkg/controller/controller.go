package controller

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/suzuki-shunsuke/tengo-tester/pkg/config"
	"github.com/suzuki-shunsuke/tengo-tester/pkg/constant"
	"github.com/suzuki-shunsuke/tengo-tester/pkg/expr"
	"github.com/suzuki-shunsuke/tengo-tester/pkg/template"
)

func (ctrl Controller) testByScript(scr string, test config.Test, entryName string, result map[string]interface{}) error {
	prog := expr.New(scr)
	result, err := prog.Run(map[string]interface{}{
		"result":     result,
		"entry_name": entryName,
		"test_name":  test.Name,
	})
	if err != nil {
		return fmt.Errorf("test a script: %w", err)
	}
	res, ok := result["err_msg"]
	if !ok {
		return nil
	}
	msg, ok := res.(string)
	if !ok {
		return fmt.Errorf(`the variable "err_msg" should be string: %+v`, res)
	}
	if msg != "" {
		return errors.New(msg)
	}
	return nil
}

func (ctrl Controller) readScript(p, wd string) (string, error) {
	if !filepath.IsAbs(p) {
		p = filepath.Join(wd, p)
	}
	f, err := ctrl.FileReader.Open(p)
	if err != nil {
		return "", fmt.Errorf("open a script "+p+": %w", err)
	}
	defer f.Close()
	a, err := ioutil.ReadAll(f)
	if err != nil {
		return "", fmt.Errorf("read a script "+p+": %w", err)
	}
	return string(a), nil
}

func (ctrl Controller) testVar(test config.Test, result map[string]interface{}) error {
	v, ok := result[test.VarName]
	if !ok {
		return errors.New(`the variable ` + test.VarName + " is undefined")
	}
	if test.Equal != nil {
		if diff := cmp.Diff(test.Equal, v); diff != "" {
			return errors.New("the value of the variable " + test.VarName + " is unexpected (-want, +got)\n" + strings.TrimSpace(diff))
		}
	}
	if test.IsNil {
		if v != nil {
			return errors.New("the variable " + test.VarName + " should be nil")
		}
	}
	if test.IsNotNil {
		if v == nil {
			return errors.New("the variable " + test.VarName + " should not be nil")
		}
	}
	return nil
}

func (ctrl Controller) test(wd string, test config.Test, entryName string, result map[string]interface{}) error {
	scr := test.Script
	if test.ScriptFile != "" {
		s, err := ctrl.readScript(test.ScriptFile, wd)
		if err != nil {
			return err
		}
		scr = s
	}
	if test.VarName != "" {
		if err := ctrl.testVar(test, result); err != nil {
			return err
		}
	}
	return ctrl.testByScript(scr, test, entryName, result)
}

func (ctrl Controller) testEntry(wd string, entry config.Entry, testErrorTpl template.Template) error {
	p := entry.ScriptFile
	if !filepath.IsAbs(p) {
		p = filepath.Join(wd, p)
	}
	f, err := ctrl.FileReader.Open(p)
	if err != nil {
		return fmt.Errorf("open a script "+p+": %w", err)
	}
	defer f.Close()
	a, err := ioutil.ReadAll(f)
	if err != nil {
		return fmt.Errorf("read a script "+p+": %w", err)
	}
	prog := expr.New(string(a))
	result, err := prog.Run(entry.Params)
	if err != nil {
		return fmt.Errorf("execute a script "+p+": %w", err)
	}
	passed := true
	for _, test := range entry.Tests {
		if err := ctrl.test(wd, test, entry.Name, result); err != nil {
			passed = false
			s, err := testErrorTpl.Render(map[string]interface{}{
				"entry_name": entry.Name,
				"test_name":  test.Name,
				"err":        err,
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, "render a task error message:", err)
				continue
			}
			fmt.Fprintln(os.Stderr, s)
		}
	}
	if !passed {
		return errors.New("")
	}
	return nil
}

func (ctrl Controller) Run(ctx context.Context, wd string) error {
	passed := true
	testErrorTpl, err := template.Compile(constant.TestError)
	if err != nil {
		return err
	}
	entryErrorTpl, err := template.Compile(constant.EntryError)
	if err != nil {
		return err
	}
	for _, entry := range ctrl.Config.Entries {
		err := ctrl.testEntry(wd, entry, testErrorTpl)
		if err != nil {
			passed = false
			s, err := entryErrorTpl.Render(map[string]interface{}{
				"entry_name": entry.Name,
				"err":        err,
			})
			if err != nil {
				fmt.Fprintln(os.Stderr, "render entry error message:", err)
				continue
			}
			fmt.Fprintln(os.Stderr, s)
		}
	}
	if !passed {
		return errors.New("")
	}
	return nil
}
