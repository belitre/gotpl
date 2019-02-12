package tpl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig"
	"github.com/ghodss/yaml"
	"k8s.io/helm/pkg/chartutil"
	"k8s.io/helm/pkg/strvals"
)

func getFunctions() template.FuncMap {
	f := sprig.TxtFuncMap()

	// from Helm!
	extra := template.FuncMap{
		"toToml":   chartutil.ToToml,
		"toYaml":   chartutil.ToYaml,
		"fromYaml": chartutil.FromYaml,
		"toJson":   chartutil.ToJson,
		"fromJson": chartutil.FromJson,
	}

	for k, v := range extra {
		f[k] = v
	}

	return f
}

func executeTemplates(values map[string]interface{}, tplFiles []string, isStrict bool) (string, error) {
	buf := bytes.NewBuffer(nil)
	for _, f := range tplFiles {
		tpl := template.New(path.Base(f)).Funcs(getFunctions())
		if isStrict {
			tpl.Option("missingkey=error")
		}

		tpl, err := tpl.ParseFiles(f)
		if err != nil {
			return "", fmt.Errorf("Error parsing template(s): %v", err)
		}

		err = tpl.Execute(buf, values)
		if err != nil {
			return "", fmt.Errorf("Failed to parse standard input: %v", err)
		}
	}

	// Work around to remove the "<no value>" go templates add.
	return strings.Replace(buf.String(), "<no value>", "", -1), nil
}

// ParseTemplate reads YAML or JSON documents from valueFiles, and extra values
// from setValues, and it uses those values for the tplFileName template,
// and writes the executed templates to the out stream.
func ParseTemplate(tplFileNames []string, valueFiles []string, setValues []string, isStrict bool) error {
	values, err := vals(valueFiles, setValues)
	if err != nil {
		return err
	}

	result, err := executeTemplates(values, tplFileNames, isStrict)
	if err != nil {
		return err
	}

	fmt.Println(result)
	return nil
}

// HELM CODE
// I really like how you can set values with helm... so using their code:
// https://github.com/kubernetes/helm/blob/master/cmd/helm/install.go

// vals merges values from files specified via -f/--values and
// directly via --set, marshaling them to YAML
func vals(valueFiles []string, values []string) (map[string]interface{}, error) {
	base := map[string]interface{}{}

	// User specified a values files via -f/--values
	for _, filePath := range valueFiles {
		currentMap := map[string]interface{}{}

		var bytes []byte
		var err error
		if strings.TrimSpace(filePath) == "-" {
			bytes, err = ioutil.ReadAll(os.Stdin)
		} else {
			bytes, err = ioutil.ReadFile(filePath)
		}

		if err != nil {
			return map[string]interface{}{}, err
		}

		if err := yaml.Unmarshal(bytes, &currentMap); err != nil {
			return map[string]interface{}{}, fmt.Errorf("failed to parse %s: %s", filePath, err)
		}
		// Merge with the previous map
		base = mergeValues(base, currentMap)
	}

	// User specified a value via --set
	for _, value := range values {
		if err := strvals.ParseInto(value, base); err != nil {
			return map[string]interface{}{}, fmt.Errorf("failed parsing --set data: %s", err)
		}
	}

	return base, nil
}

// Merges source and destination map, preferring values from the source map
func mergeValues(dest map[string]interface{}, src map[string]interface{}) map[string]interface{} {
	for k, v := range src {
		// If the key doesn't exist already, then just set the key to that value
		if _, exists := dest[k]; !exists {
			dest[k] = v
			continue
		}
		nextMap, ok := v.(map[string]interface{})
		// If it isn't another map, overwrite the value
		if !ok {
			dest[k] = v
			continue
		}
		// Edge case: If the key exists in the destination, but isn't a map
		destMap, isMap := dest[k].(map[string]interface{})
		// If the source map has a map for this key, prefer it
		if !isMap {
			dest[k] = v
			continue
		}
		// If we got to this point, it is a map in both, so merge them
		dest[k] = mergeValues(destMap, nextMap)
	}
	return dest
}
