package api

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type SwaggerDoc struct {
	Paths map[string]map[string]Operation `json:"paths"`
}

type Operation struct {
	Summary     string              `json:"summary"`
	Description string              `json:"description"`
	Parameters  []Parameter         `json:"parameters"`
	Responses   map[string]Response `json:"responses"`
}

type Parameter struct {
	Name        string `json:"name"`
	In          string `json:"in"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Type        string `json:"type"`
}

type Response struct {
	Description string `json:"description"`
}

func LoadSwaggerDoc(filename string) (*SwaggerDoc, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open swagger file: %w", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read swagger file: %w", err)
	}

	var swagger SwaggerDoc
	err = json.Unmarshal(data, &swagger)
	if err != nil {
		return nil, fmt.Errorf("failed to parse swagger file: %w", err)
	}

	return &swagger, nil
}
