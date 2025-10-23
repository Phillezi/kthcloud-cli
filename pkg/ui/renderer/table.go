package renderer

import (
	"fmt"
	"io"
	"reflect"
	"text/tabwriter"
)

func renderTable(w io.Writer, obj any) error {
	switch v := obj.(type) {
	case []DeploymentLike:
		return renderDeployments(w, v)
	default:
		return renderGeneric(w, obj)
	}
}

func renderGeneric(w io.Writer, obj any) error {
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Pointer {
		v = v.Elem()
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		return renderSlice(w, v)
	case reflect.Struct:
		return renderStruct(w, v)
	default:
		_, err := fmt.Fprintln(w, derefValue(v))
		return err
	}
}

func renderSlice(w io.Writer, v reflect.Value) error {
	if v.Len() == 0 {
		fmt.Fprintln(w, "No results found.")
		return nil
	}

	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	elem := v.Index(0)
	if elem.Kind() == reflect.Pointer {
		elem = elem.Elem()
	}

	if elem.Kind() != reflect.Struct {
		for i := 0; i < v.Len(); i++ {
			fmt.Fprintln(tw, derefValue(v.Index(i)))
		}
		return tw.Flush()
	}

	t := elem.Type()
	cols := make([]string, 0, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		cols = append(cols, t.Field(i).Name)
	}
	fmt.Fprintln(tw, joinTabs(cols...))

	for i := 0; i < v.Len(); i++ {
		row := v.Index(i)
		if row.Kind() == reflect.Pointer {
			row = row.Elem()
		}
		values := make([]string, 0, t.NumField())
		for j := 0; j < t.NumField(); j++ {
			values = append(values, fmt.Sprintf("%v", derefValue(row.Field(j))))
		}
		fmt.Fprintln(tw, joinTabs(values...))
	}

	return tw.Flush()
}

func renderStruct(w io.Writer, v reflect.Value) error {
	t := v.Type()
	tw := tabwriter.NewWriter(w, 0, 0, 2, ' ', 0)
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		fmt.Fprintf(tw, "%s:\t%v\n", f.Name, derefValue(v.Field(i)))
	}
	return tw.Flush()
}

func derefValue(v reflect.Value) any {
	for v.Kind() == reflect.Pointer {
		if v.IsNil() {
			return nil
		}
		v = v.Elem()
	}
	return v.Interface()
}

func joinTabs(fields ...string) string {
	out := ""
	for i, f := range fields {
		if i > 0 {
			out += "\t"
		}
		out += f
	}
	return out
}
