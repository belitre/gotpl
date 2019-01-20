package tpl

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTemplate(t *testing.T) {
	type io struct {
		Input     []string
		Templates []string
		Output    string
		Strict    bool
		Error     bool
	}

	tests := []io{
		{
			Input:     []string{"test=value"},
			Templates: []string{"{{.test}}"},
			Output:    "value",
			Strict:    false,
			Error:     false,
		},
		{
			Input:     []string{},
			Templates: []string{`{{.foo | default "bar"}}`},
			Output:    "bar",
			Strict:    false,
			Error:     false,
		},
		{
			Input:     []string{"foo="},
			Templates: []string{`{{.foo | default "bar"}}`},
			Output:    "bar",
			Strict:    false,
			Error:     false,
		},
		{
			Input:     []string{"foo="},
			Templates: []string{`test {{.foo}}`},
			Output:    "test ",
			Strict:    false,
			Error:     false,
		},
		{
			Input:     []string{"foo=bleh"},
			Templates: []string{`{{.foo | default "bar"}}`},
			Output:    "bleh",
			Strict:    false,
			Error:     false,
		},
		{
			Input:     []string{"user=u", "password=p"},
			Templates: []string{`{{ (printf "%s:%s" .user .password) | b64enc }}`},
			Output:    "dTpw",
			Strict:    false,
			Error:     false,
		},
		{
			Input:     []string{"name=Max", "age=15"},
			Templates: []string{"Hello {{.name}}, of {{.age}} years old"},
			Output:    "Hello Max, of 15 years old",
			Strict:    false,
			Error:     false,
		},
		{
			Input:     []string{"legumes={potato,onion,cabbage}"},
			Templates: []string{"Legumes:{{ range $index, $el := .legumes}}{{if $index}},{{end}} {{$el}}{{end}}"},
			Output:    "Legumes: potato, onion, cabbage",
			Strict:    false,
			Error:     false,
		},
		{
			Input:     []string{"test=value"},
			Templates: []string{"{{.test}}"},
			Output:    "value",
			Strict:    true,
			Error:     false,
		},
		{
			Input:     []string{},
			Templates: []string{`{{.foo | default "bar"}}`},
			Output:    "",
			Strict:    true,
			Error:     true,
		},
		{
			Input:     []string{"foo="},
			Templates: []string{`{{.foo | default "bar"}}`},
			Output:    "bar",
			Strict:    true,
			Error:     false,
		},
		{
			Input:     []string{"foo="},
			Templates: []string{`test {{.foo}}`},
			Output:    "test ",
			Strict:    true,
			Error:     false,
		},
		{
			Input:     []string{"foo=bleh"},
			Templates: []string{`{{.foo | default "bar"}}`},
			Output:    "bleh",
			Strict:    true,
			Error:     false,
		},
		{
			Input:     []string{"user=u", "password=p"},
			Templates: []string{`{{ (printf "%s:%s" .user .password) | b64enc }}`},
			Output:    "dTpw",
			Strict:    true,
			Error:     false,
		},
		{
			Input:     []string{"name=Max", "age=15"},
			Templates: []string{"Hello {{.name}}, of {{.age}} years old"},
			Output:    "Hello Max, of 15 years old",
			Strict:    true,
			Error:     false,
		},
		{
			Input:     []string{"legumes={potato,onion,cabbage}"},
			Templates: []string{"Legumes:{{ range $index, $el := .legumes}}{{if $index}},{{end}} {{$el}}{{end}}"},
			Output:    "Legumes: potato, onion, cabbage",
			Strict:    true,
			Error:     false,
		},
		{
			Input: []string{"user=myuser", "password=mypass"},
			Templates: []string{
				"{{ .user }}\n",
				"{{ .password | b64enc }}",
			},
			Output: "myuser\nbXlwYXNz",
			Strict: false,
			Error:  false,
		},
	}

	for _, test := range tests {
		fileNames := []string{}
		for _, templ := range test.Templates {
			tplFile, err := ioutil.TempFile("", "")
			assert.Nil(t, err)
			defer func() {
				os.Remove(tplFile.Name())
			}()
			_, err = tplFile.WriteString(templ)
			assert.Nil(t, err)
			fileNames = append(fileNames, tplFile.Name())
			tplFile.Close()
		}
		values, err := vals(nil, test.Input)
		assert.Nil(t, err)
		output, err := executeTemplates(values, fileNames, test.Strict)
		if test.Error {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.Output, output)
		}
	}
}
