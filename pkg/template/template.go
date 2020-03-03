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

func ExecuteTemplate(template string, params map[string]string) string {
	output := template

	for key, value := range params {
		output = strings.Replace(output, "${"+key+"}", value, -1)
	}

	return output
}
