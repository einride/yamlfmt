package main

import (
	"testing"

	"gotest.tools/v3/assert"
)

func TestFormat(t *testing.T) {
	for _, tt := range []struct {
		name   string
		indent int
		input  string
		output string
		err    string
	}{
		{
			name: "empty",
		},
		{
			name:   "one row with linux newline",
			indent: 2,
			input: `name: foo
`,
			output: `name: foo
`,
		},
		{
			name:   "one row without newline",
			indent: 2,
			input:  "name: foo",
			output: `name: foo
`,
		},
		{
			name:   "two row with linux newline",
			indent: 2,
			input: `name: foo
age: 33
`,
			output: `name: foo
age: 33
`,
		},
		{
			name:   "two row with linux newline between",
			indent: 2,
			input: `
name: foo

age: 33
`,
			output: `name: foo

age: 33
`,
		},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			o, err := format([]byte(tt.input), tt.indent)
			if tt.err != "" {
				assert.ErrorContains(t, err, tt.err)
			} else {
				assert.DeepEqual(t, string(o), tt.output)
			}
		})
	}
}
