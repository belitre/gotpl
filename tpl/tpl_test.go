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
	}

	tests := []io{
		{
			Input:    []string{"test=value"},
			Template: "{{.test}}",
			Output:   "value",
		},
		{
			Input:    []string{},
			Template: `{{.foo | default "bar"}}`,
			Output:   "bar",
		},
		{
			Input:    []string{"foo="},
			Template: `{{.foo | default "bar"}}`,
			Output:   "bar",
		},
		{
			Input:    []string{"foo=bleh"},
			Template: `{{.foo | default "bar"}}`,
			Output:   "bleh",
		},
		{
			Input:    []string{"user=u", "password=p"},
			Template: `{{ (printf "%s:%s" .user .password) | b64enc }}`,
			Output:   "dTpw",
		},
		{
			Input:    []string{"name=Max", "age=15"},
			Template: "Hello {{.name}}, of {{.age}} years old",
			Output:   "Hello Max, of 15 years old",
		},
		{
			Input:    []string{"legumes={potato,onion,cabbage}"},
			Template: "Legumes:{{ range $index, $el := .legumes}}{{if $index}},{{end}} {{$el}}{{end}}",
			Output:   "Legumes: potato, onion, cabbage",
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
		output, err := executeTemplates(values, tplFile.Name())
		assert.Nil(t, err)
		assert.Equal(t, test.Output, output)
	}
}
