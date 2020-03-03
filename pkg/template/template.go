package template

import (
	"fmt"
	"io/ioutil"
	"strings"
)

func LoadTemplate(path string) string {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Println("Err: ", err)
	}

	return string(data)
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
