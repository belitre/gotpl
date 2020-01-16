package tpl

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/belitre/gotpl/commands/options"
	"github.com/ghodss/yaml"
	"github.com/otiai10/copy"
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

func executeSingleTemplate(values map[string]interface{}, tplFile string, isStrict bool) (string, error) {
	buf := bytes.NewBuffer(nil)
	tpl := template.New(path.Base(tplFile)).Funcs(getFunctions())
	if isStrict {
		tpl.Option("missingkey=error")
	}

	tpl, err := tpl.ParseFiles(tplFile)
	if err != nil {
		return "", fmt.Errorf("Error parsing template(s): %v", err)
	}

	err = tpl.Execute(buf, values)
	if err != nil {
		return "", fmt.Errorf("Failed to parse standard input: %v", err)
	}

	// Work around to remove the "<no value>" go templates add.
	return strings.Replace(buf.String(), "<no value>", "", -1), nil
}

func executeTemplates(values map[string]interface{}, tplFileNames []string, isStrict bool, outputPath string) (string, error) {
	var result string
	var tmpDir string
	var err error
	if len(outputPath) > 0 {
		tmpDir, err = ioutil.TempDir("", "")

		if err != nil {
			return "", err
		}
		defer func() {
			os.RemoveAll(tmpDir)
		}()
	}

	listFiles, err := getListFiles(tplFileNames)

	if err != nil {
		return "", err
	}

	for _, f := range listFiles {
		r, err := executeSingleTemplate(values, f.src, isStrict)
		if err != nil {
			return "", err
		}
		if len(tmpDir) > 0 {
			err = saveFile(path.Join(tmpDir, f.dest), r, f.perm)
			if err != nil {
				return "", err
			}
		} else {
			if len(result) == 0 {
				result = r
			} else {
				result = fmt.Sprintf("%s\n%s", result, r)
			}
		}
	}

	if len(tmpDir) > 0 {
		err := copy.Copy(tmpDir, outputPath)
		if err != nil {
			return "", err
		}
	}

	return result, nil
}

func saveFile(path string, contents string, perm os.FileMode) error {
	if err := os.MkdirAll(filepath.Dir(path), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}

	defer f.Close()

	if err = f.Chmod(perm); err != nil {
		return err
	}

	if _, err = f.WriteString(contents); err != nil {
		return err
	}

	return nil
}

type SrcDest struct {
	src  string
	dest string
	perm os.FileMode
}

// ParseTemplate reads YAML or JSON documents from valueFiles, and extra values
// from setValues, and it uses those values for the tplFileName template,
// and writes the executed templates to the out stream.
func ParseTemplate(tplFileNames []string, opts *options.Options) error {
	values, err := vals(opts.ValueFiles, opts.SetValues)
	if err != nil {
		return err
	}

	result, err := executeTemplates(values, tplFileNames, opts.IsStrict, opts.OutputPath)
	if err != nil {
		return err
	}

	if len(result) > 0 {
		fmt.Println(result)
	}

	return nil
}

func getListFiles(tplFileNames []string) ([]*SrcDest, error) {
	listFiles := []*SrcDest{}

	for _, f := range tplFileNames {
		cleanPath := path.Clean(f)
		info, err := os.Stat(cleanPath)
		if err != nil {
			return nil, err
		}

		if info.IsDir() {
			err := filepath.Walk(f,
				func(p string, i os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !i.IsDir() {
						srcdst := &SrcDest{
							src:  p,
							dest: strings.Replace(p, cleanPath, "", 1),
							perm: i.Mode(),
						}
						listFiles = append(listFiles, srcdst)
					}
					return nil
				})
			if err != nil {
				return nil, err
			}

		} else {
			srcdst := &SrcDest{
				src:  cleanPath,
				dest: filepath.Base(cleanPath),
				perm: info.Mode(),
			}
			listFiles = append(listFiles, srcdst)
		}
	}

	return listFiles, nil
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
