# tengo-tester

CLI to test a Tengo script

## Overview

tengo-tester is a CLI tool to test [Tengo](https://github.com/d5/tengo) scripts.
Tengo is a fast script language for Go.

We can define tests of Tengo scripts in the `tengo-tester` configuration file and test scripts by `tengo-tester run`.

## Limitation

https://github.com/d5/tengo/blob/master/docs/objects.md#user-object-types

`tengo-tester` doesn't support the Tengo's User Extension.

On the other hand, all standard library can be used.

## Install

Download from [GitHub Releases](https://github.com/suzuki-shunsuke/tengo-tester/releases)

```
$ tengo-tester --version
tengo-tester version 0.1.0
```

## Getting Started

Write a Tengo script which we want to test.

foo.tengo

```
name := "foo"
```

Generate the configuration file `.tengo-tester.yaml` by `tengo-tester init`.

```
$ tengo-tester init
```

Edit the configuration file.

.tengo-tester.yaml

```yaml
---
entries:
- name: main
  script_file: foo.tengo
  tests:
  - name: test foo
    var_name: name
    equal: foo
```

Run test.

```
$ tengo-tester run
```

Change the configuraiton `equal` from `foo` to `fo` for the test to fail.

```
$ tengo-tester run
test fails
entry_name: main
test_name: test
the value of the variable foo is unexpected (-want, +got)
string(
-       "fo",
+       "foo",
  )
entry fails
entry_name: main

FATA[0000]
exit status 1
```

## Configuration Reference

```yaml
---
entries:
- name: main # entry name
  script_file: foo.tengo # the tested script file path
  params: # the parameter of the tengo script.
    var_name: var_value
  tests:
  - name: foo # test name
    script: |
      # tengo script to test the result
      # If the test fails, set the error message to the variable "err_msg".
      fmt := import("fmt")
      err_msg := ""
      if result.foo != "foo" {
        err_msg = fmt.sprintf("foo = %v, wanted foo", result.foo)
      }
    script_file: foo_test.tengo # file path to a test script
    var_name: foo # the variable name to be tested
    equal: foo # When `equal` isn't nil, the test fails if the variable isn't equal to the value of `equal`.
    is_nil: true # When `is_nil` is true, the test fails if the variable isn't nil.
    is_not_nil: true # When `is_not_nil` is true, the test fails if the variable is nil.
```

## Configuration file path

The configuration file path can be specified with the `--config (-c)` option.
If the confgiuration file path isn't specified, the file named `.tengo-tester.yml` or `.tengo-tester.yaml` would be searched from the current directory to the root directory.

## the base directory of the relative path

There are some configuration which are file paths.
If a file path is a relative path, the base directory of the relative path is the directory where the configuration file exists.

## LICENSE

[MIT](LICENSE)
