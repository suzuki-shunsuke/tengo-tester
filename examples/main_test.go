package examples_test

import (
	"testing"

	"gotest.tools/v3/icmd"
)

func TestBuildflow(t *testing.T) { //nolint:funlen
	data := []struct {
		title string
		file  string
		exp   icmd.Expected
	}{
		{
			title: "hello world",
			file:  "hello_world.yaml",
		},
	}
	for _, d := range data {
		d := d
		t.Run(d.title, func(t *testing.T) {
			result := icmd.RunCmd(icmd.Command("tengo-tester", "run", "-c", d.file))
			result.Assert(t, d.exp)
		})
	}
}
