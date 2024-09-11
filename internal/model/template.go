package model

import (
	"bytes"
	"text/template"
)

type TemplateData struct {
	BaseURL  string
	Realm    string
	ClientID string
}

func NewTemplate(baseURL string,
	realm string,
	clientID string) *TemplateData {
	return &TemplateData{
		BaseURL:  baseURL,
		Realm:    realm,
		ClientID: clientID,
	}
}

func (t *TemplateData) Replace(templateStr string) (string, error) {

	tmpl, err := template.New("script").Parse(templateStr)
	if err != nil {
		return "", err
	}

	var outputBuffer bytes.Buffer
	err = tmpl.Execute(&outputBuffer, t)
	if err != nil {
		return "", err
	}

	return outputBuffer.String(), nil
}
