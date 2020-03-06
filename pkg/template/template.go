package template

import (
	"io/ioutil"
	"strings"
)

func LoadTemplate(path string) (string, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func ExecuteTemplate(template string, vars map[string]string, params map[string]string) string {
	output := template

	for key, value := range vars {
		output = strings.Replace(output, "${global."+key+"}", value, -1)
	}

	for key, value := range params {
		output = strings.Replace(output, "${params."+key+"}", value, -1)
	}

	return output
}
