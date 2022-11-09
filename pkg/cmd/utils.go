package cmd

import (
	"encoding/csv"
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

// Table 格式化方式
const (
	Table = "table"
	Yaml  = "yaml"
	Json  = "json"
	Csv   = "csv"
)

// RenderData 将数据用指定形式输出
func RenderData(format string, header []string, data [][]string, raw any, writer io.Writer) error {
	switch format {
	case Yaml:
		yamlEncoder := yaml.NewEncoder(writer)
		err := yamlEncoder.Encode(raw)
		if err != nil {
			return err
		}
	case Csv:
		csvWriter := csv.NewWriter(writer)
		err := csvWriter.WriteAll(data)
		if err != nil {
			return err
		}
		err = csvWriter.Error()
		if err != nil {
			return err
		}
	case Json:
		jsonEncoder := jsoniter.NewEncoder(writer)
		jsonEncoder.SetIndent("", "  ")
		err := jsonEncoder.Encode(raw)
		if err != nil {
			return err
		}
	case Table:
		fallthrough
	default:
		table := tablewriter.NewWriter(writer)
		table.SetAutoWrapText(false)
		table.SetAlignment(tablewriter.ALIGN_LEFT)
		table.AppendBulk(data)
		table.SetHeader(header)
		table.SetRowLine(true)
		table.Render()
	}

	return nil
}

func NoArgs(cmd *cobra.Command, args []string) error {
	if err := cobra.NoArgs(cmd, args); err != nil {
		_ = cmd.Usage()
		return err
	}

	return nil
}

func ExactArgs(n int) cobra.PositionalArgs {
	return func(cmd *cobra.Command, args []string) error {
		if err := cobra.ExactArgs(n)(cmd, args); err != nil {
			_ = cmd.Usage()
			return err
		}
		return nil
	}
}
