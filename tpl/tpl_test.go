package tpl

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestTemplate(t *testing.T) {
	type io struct {
		Input    []string
		Template string
		Output   string
		Strict   bool
		Error    bool
	}

	tests := []io{
		{
			Input:    []string{"test=value"},
			Template: "{{.test}}",
			Output:   "value",
			Strict:   false,
			Error:    false,
		},
		{
			Input:    []string{},
			Template: `{{.foo | default "bar"}}`,
			Output:   "bar",
			Strict:   false,
			Error:    false,
		},
		{
			Input:    []string{"foo="},
			Template: `{{.foo | default "bar"}}`,
			Output:   "bar",
			Strict:   false,
			Error:    false,
		},
		{
			Input:    []string{"foo="},
			Template: `test {{.foo}}`,
			Output:   "test ",
			Strict:   false,
			Error:    false,
		},
		{
			Input:    []string{"foo=bleh"},
			Template: `{{.foo | default "bar"}}`,
			Output:   "bleh",
			Strict:   false,
			Error:    false,
		},
		{
			Input:    []string{"user=u", "password=p"},
			Template: `{{ (printf "%s:%s" .user .password) | b64enc }}`,
			Output:   "dTpw",
			Strict:   false,
			Error:    false,
		},
		{
			Input:    []string{"name=Max", "age=15"},
			Template: "Hello {{.name}}, of {{.age}} years old",
			Output:   "Hello Max, of 15 years old",
			Strict:   false,
			Error:    false,
		},
		{
			Input:    []string{"legumes={potato,onion,cabbage}"},
			Template: "Legumes:{{ range $index, $el := .legumes}}{{if $index}},{{end}} {{$el}}{{end}}",
			Output:   "Legumes: potato, onion, cabbage",
			Strict:   false,
			Error:    false,
		},
		{
			Input:    []string{"test=value"},
			Template: "{{.test}}",
			Output:   "value",
			Strict:   true,
			Error:    false,
		},
		{
			Input:    []string{},
			Template: `{{.foo | default "bar"}}`,
			Output:   "",
			Strict:   true,
			Error:    true,
		},
		{
			Input:    []string{"foo="},
			Template: `{{.foo | default "bar"}}`,
			Output:   "bar",
			Strict:   true,
			Error:    false,
		},
		{
			Input:    []string{"foo="},
			Template: `test {{.foo}}`,
			Output:   "test ",
			Strict:   true,
			Error:    false,
		},
		{
			Input:    []string{"foo=bleh"},
			Template: `{{.foo | default "bar"}}`,
			Output:   "bleh",
			Strict:   true,
			Error:    false,
		},
		{
			Input:    []string{"user=u", "password=p"},
			Template: `{{ (printf "%s:%s" .user .password) | b64enc }}`,
			Output:   "dTpw",
			Strict:   true,
			Error:    false,
		},
		{
			Input:    []string{"name=Max", "age=15"},
			Template: "Hello {{.name}}, of {{.age}} years old",
			Output:   "Hello Max, of 15 years old",
			Strict:   true,
			Error:    false,
		},
		{
			Input:    []string{"legumes={potato,onion,cabbage}"},
			Template: "Legumes:{{ range $index, $el := .legumes}}{{if $index}},{{end}} {{$el}}{{end}}",
			Output:   "Legumes: potato, onion, cabbage",
			Strict:   true,
			Error:    false,
		},
	}

	for _, test := range tests {
		tplFile, err := ioutil.TempFile("", "")
		assert.Nil(t, err)
		defer func() { os.Remove(tplFile.Name()) }()

		_, err = tplFile.WriteString(test.Template)
		assert.Nil(t, err)
		tplFile.Close()

		values, err := vals(nil, test.Input)
		assert.Nil(t, err)
		output, err := executeTemplates(values, tplFile.Name(), test.Strict)
		if test.Error {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.Output, output)
		}
	}
}
