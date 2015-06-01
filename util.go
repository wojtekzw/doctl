package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v1"
)

func WriteCollectionOutput(data interface{}, headers []string, format string, row func(interface{}) []interface{}) {
	switch OutputFormat {
	case "table":
		t := tableOutput{
			headers:   headers,
			rowFormat: format,
			rowValues: row,
			w:         tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0),
		}

		switch dataArray := data.(type) {
		case []interface{}:
			data
		default:

		}

		defer t.flush()
		t.header()
		for _, datum := range dataArray {
			t.writeln(datum)
		}
	case "json":
		output, err := json.Marshal(data)
		if err != nil {
			fmt.Printf("JSON Encoding Error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("%s", string(output))
	case "yaml":
		output, err := yaml.Marshal(data)
		if err != nil {
			fmt.Printf("YAML Encoding Error: %s", err)
			os.Exit(1)
		}
		fmt.Printf("%s", string(output))
	}
}

type tableOutput struct {
	headers   []string
	rowFormat string
	rowValues func(interface{}) []interface{}
	w         *tabwriter.Writer
}

func (t *tableOutput) header() {
	fmt.Fprintln(t.w, strings.Join(t.headers, "\t"))
}

func (t *tableOutput) writeln(datum interface{}) {
	fields := t.rowValues(datum)
	fmt.Fprintf(t.w, t.rowFormat, fields...)
}

func (t *tableOutput) flush() {
	t.w.Flush()
}

//
// Legacy printers
//

type CLIOutput struct {
	w *tabwriter.Writer
}

func NewCLIOutput() *CLIOutput {
	return &CLIOutput{
		w: tabwriter.NewWriter(os.Stdout, 0, 8, 2, '\t', 0),
	}
}

func WriteOutput(data interface{}) {
	var output []byte
	var err error

	switch OutputFormat {
	case "json":
		output, err = json.Marshal(data)
		if err != nil {
			fmt.Printf("JSON Encoding Error: %s", err)
			os.Exit(1)
		}

	case "yaml":
		output, err = yaml.Marshal(data)
		if err != nil {
			fmt.Printf("YAML Encoding Error: %s", err)
			os.Exit(1)
		}
	}
	fmt.Printf("%s", string(output))
}

func (c *CLIOutput) Header(a ...string) {
	fmt.Fprintln(c.w, strings.Join(a, "\t"))
}

func (c *CLIOutput) Writeln(format string, a ...interface{}) {
	fmt.Fprintf(c.w, format, a...)
}

func (c *CLIOutput) Flush() {
	c.w.Flush()
}
