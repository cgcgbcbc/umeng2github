package main

import (
	"bytes"
	"text/template"
)

const DefaultBodyTemplate = `
发生次数 | 首次发生时间 | 版本
-----|-----|-----
{{.Count}} | {{.FirstDate.Format "2006-01-02 15:04:05"}} | {{.Version}}

StackTrace:
` +
	"```\n{{.StackTrace}}\n```"

func FormatRecord(record Record) (result string, err error) {
	tmpl, err := template.New("").Parse(DefaultBodyTemplate)
	if err != nil {
		return
	}

	var b bytes.Buffer
	err = tmpl.Execute(&b, record)
	if err != nil {
		return
	}

	result = b.String()

	return
}
