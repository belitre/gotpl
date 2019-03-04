package tpl

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/belitre/gotpl/commands/options"
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
				"{{ .user }}",
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
		output, err := executeTemplates(values, fileNames, test.Strict, "")
		if test.Error {
			assert.NotNil(t, err)
		} else {
			assert.Nil(t, err)
			assert.Equal(t, test.Output, output)
		}
	}
}

func TestOutputFolder(t *testing.T) {
	tplFiles := []string{
		"test/tpl1.tpl",
		"test/tpl2.tpl",
		"test/folder",
	}

	outPath := "result"

	err := os.Mkdir(outPath, os.ModePerm)
	assert.NoError(t, err)

	defer os.RemoveAll(outPath)

	opts := &options.Options{
		IsStrict:   true,
		OutputPath: outPath,
		SetValues: []string{
			"name=paco",
			"test=cat",
			"bleh=pato",
		},
		ValueFiles: []string{},
	}

	err = ParseTemplate(tplFiles, opts)
	assert.NoError(t, err)

	type result struct {
		fileName string
		content  string
	}

	results := []result{
		result{
			fileName: "result/tpl1.tpl",
			content:  "Hello paco",
		},
		result{
			fileName: "result/tpl2.tpl",
			content:  "Bye paco",
		},
		result{
			fileName: "result/tpl3.tpl",
			content:  "Hello cat",
		},
		result{
			fileName: "result/tpl4.tpl",
			content:  "Bye cat",
		},
		result{
			fileName: "result/sub/tpl5.tpl",
			content:  "Hello pato",
		},
		result{
			fileName: "result/sub/tpl6.tpl",
			content:  "Bye pato",
		},
	}

	for _, test := range results {
		c, err := ioutil.ReadFile(test.fileName)
		assert.NoError(t, err)
		assert.Equal(t, test.content, string(c[:]))
	}
}
