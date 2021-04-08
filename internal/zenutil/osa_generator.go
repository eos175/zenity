// +build tools

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/dchest/jsmin"
)

func main() {
	dir := os.Args[1]

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}

	var str strings.Builder

	for _, file := range files {
		name := file.Name()

		str.WriteString("\n" + `{{define "`)
		str.WriteString(strings.TrimSuffix(name, filepath.Ext(name)))
		str.WriteString(`" -}}` + "\n")

		data, err := ioutil.ReadFile(filepath.Join(dir, name))
		if err != nil {
			log.Fatal(err)
		}
		data, err = minify(data)
		if err != nil {
			log.Fatal(err)
		}

		str.Write(data)
		str.WriteString("\n{{- end}}")
	}

	out, err := os.Create("osa_generated.go")
	if err != nil {
		log.Fatal(err)
	}

	err = generator.Execute(out, str.String())
	if err != nil {
		log.Fatal(err)
	}

	err = out.Close()
	if err != nil {
		log.Fatal(err)
	}
}

func minify(data []byte) ([]byte, error) {
	var templates [][]byte
	var buf []byte

	for {
		i := bytes.Index(data, []byte("{{"))
		if i < 0 {
			break
		}
		j := bytes.Index(data[i+len("{{"):], []byte("}}"))
		if j < 0 {
			break
		}
		templates = append(templates, data[i:i+j+len("{{}}")])
		buf = append(buf, data[:i]...)
		buf = append(buf, []byte("TEMPLATE")...)
		data = data[i+j+len("{{}}"):]
	}
	buf = append(buf, data...)

	buf, err := jsmin.Minify(buf)
	if err != nil {
		return nil, err
	}

	var res []byte
	for _, t := range templates {
		i := bytes.Index(buf, []byte("TEMPLATE"))
		res = append(res, buf[:i]...)
		res = append(res, t...)
		buf = buf[i+len("TEMPLATE"):]
	}
	return append(res, buf...), nil
}

var generator = template.Must(template.New("").Parse(`// Code generated by zenity; DO NOT EDIT.
// +build darwin

package zenutil

import (
	"encoding/json"
	"text/template"
)

var scripts = template.Must(template.New("").Funcs(template.FuncMap{"json": func(v interface{}) (string, error) {
	b, err := json.Marshal(v)
	return string(b), err
}}).Parse(` + "`{{.}}`" + `))
`))
