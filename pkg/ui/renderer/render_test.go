package renderer

import (
	"bytes"
	"strings"
	"testing"
)

func TestRenderer_JSON(t *testing.T) {
	r := New()
	data := []DeploymentAdapter{
		{"d1", "API", "philip", "resourceRunning"},
		{"d2", "DB", "philip", "resourcePending"},
	}

	buf := &bytes.Buffer{}
	err := r.Render(data, WithOutput(Output_JSON), WithWriter(buf))
	if err != nil {
		t.Fatalf("Render JSON failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, `"ID": "d1"`) || !strings.Contains(output, `"Status": "resourceRunning"`) {
		t.Errorf("JSON output incorrect: %s", output)
	}
}

func TestRenderer_YAML(t *testing.T) {
	r := New()
	data := []DeploymentAdapter{
		{"d1", "API", "philip", "resourceRunning"},
	}

	buf := &bytes.Buffer{}
	err := r.Render(data, WithOutput(Output_YAML), WithWriter(buf))
	if err != nil {
		t.Fatalf("Render YAML failed: %v", err)
	}

	output := buf.String()
	// yaml.v3 lowercases field names by default
	if !strings.Contains(output, "id: d1") || !strings.Contains(output, "status: resourceRunning") {
		t.Errorf("YAML output incorrect: %s", output)
	}
}

func TestRenderer_Table_KnownType(t *testing.T) {
	r := New()
	data := []DeploymentAdapter{
		{"d1", "API", "philip", "resourceRunning"},
	}

	buf := &bytes.Buffer{}
	err := r.Render(data, WithOutput(Output_Table), WithWriter(buf))
	if err != nil {
		t.Fatalf("Render Table failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "ID") || !strings.Contains(output, "Status") {
		t.Errorf("Table header missing: %s", output)
	}
	if !strings.Contains(output, "d1") || !strings.Contains(output, "resourceRunning") {
		t.Errorf("Table row missing: %s", output)
	}
}

func TestRenderer_Reflection_SliceStruct(t *testing.T) {
	r := New()
	type Generic struct {
		Foo string
		Bar int
	}
	data := []Generic{
		{"hello", 42},
		{"world", 7},
	}

	buf := &bytes.Buffer{}
	err := r.Render(data, WithOutput(Output_Table), WithWriter(buf))
	if err != nil {
		t.Fatalf("Render generic slice failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Foo") || !strings.Contains(output, "Bar") {
		t.Errorf("Generic table header missing: %s", output)
	}
	if !strings.Contains(output, "hello") || !strings.Contains(output, "42") {
		t.Errorf("Generic table row missing: %s", output)
	}
	if !strings.Contains(output, "world") || !strings.Contains(output, "7") {
		t.Errorf("Generic table row missing: %s", output)
	}
}

func TestRenderer_Reflection_Struct(t *testing.T) {
	r := New()
	type Generic struct {
		Name string
		Age  int
	}
	data := Generic{"Alice", 30}

	buf := &bytes.Buffer{}
	err := r.Render(data, WithOutput(Output_Table), WithWriter(buf))
	if err != nil {
		t.Fatalf("Render generic struct failed: %v", err)
	}

	output := buf.String()
	if !strings.Contains(output, "Name:") || !strings.Contains(output, "Alice") {
		t.Errorf("Generic struct output missing Name: %s", output)
	}
	if !strings.Contains(output, "Age:") || !strings.Contains(output, "30") {
		t.Errorf("Generic struct output missing Age: %s", output)
	}
}
