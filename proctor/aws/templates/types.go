package templates

import "encoding/json"

type Template struct {
	AWSTemplateFormatVersion string
	Description              string               `json:",omitempty"`
	Parameters               map[string]Parameter `json:",omitempty"`
	Resources                map[string]Resource
	Outputs                  map[string]Output `json:",omitempty"`
}

func NewTemplate() *Template {
	return &Template{
		AWSTemplateFormatVersion: "2010-09-09",
		Parameters:               map[string]Parameter{},
		Resources:                map[string]Resource{},
		Outputs:                  map[string]Output{},
	}
}

type Parameter struct {
	Description           string
	Type                  string
	MinLength             string `json:",omitempty"`
	MaxLength             string `json:",omitempty"`
	Default               string `json:",omitempty"`
	AllowedPattern        string `json:",omitempty"`
	ConstraintDescription string `json:",omitempty"`
}

type Resource struct {
	Type       string
	Properties map[string]interface{}
}

type Output struct {
}

func (t *Template) String() string {
	bytes, err := json.Marshal(t)
	if err != nil {
		panic(err)
	}
	return string(bytes)
}

type Ref struct {
	Ref string
}

func Fn(name string, value interface{}) interface{} {
	return map[string]interface{}{
		"Fn::" + name: value,
	}
}

func FnJoin(sep string, elements ...interface{}) interface{} {
	return Fn("Join", []interface{}{sep, elements})
}
