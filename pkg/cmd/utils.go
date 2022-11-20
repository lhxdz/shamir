package cmd

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
	"golang.org/x/sys/unix"
	"gopkg.in/yaml.v2"

	"shamir/pkg/utils/code"
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

type bufferReadWriteClose struct {
	*bytes.Buffer
}

func (b *bufferReadWriteClose) Close() error {
	return nil
}

func NewBuffer(b *bytes.Buffer) *bufferReadWriteClose {
	return &bufferReadWriteClose{
		Buffer: b,
	}
}

type keyReadWriter struct {
	x io.ReadWriter
	y io.ReadWriter
}

func NewKeyReadWriter(x, y io.ReadWriter) *keyReadWriter {
	return &keyReadWriter{
		x: x,
		y: y,
	}
}

func (k *keyReadWriter) ToXYKeyDecoder() *xyKeyDecoder {
	return &xyKeyDecoder{
		x: code.NewKeyDecoder(k.x),
		y: code.NewKeyDecoder(k.y),
	}
}

func (k *keyReadWriter) toString() (string, string, error) {
	xData, err := io.ReadAll(k.x)
	if err != nil {
		return "", "", err
	}
	yData, err := io.ReadAll(k.y)
	if err != nil {
		return "", "", err
	}

	return string(xData), string(yData), nil
}

type xyKeyDecoder struct {
	x *code.KeyDecoder
	y *code.KeyDecoder
}

func (xy *xyKeyDecoder) decoder(key *code.Key) error {
	if key == nil {
		return fmt.Errorf("invalid nil point of key")
	}
	if key.X == nil || key.Y == nil {
		return fmt.Errorf("invalid nil point of x or y")
	}

	err := xy.x.Write(key.X)
	if err != nil {
		return err
	}
	err = xy.y.Write(key.Y)
	if err != nil {
		return err
	}
	return nil
}

// TaskIndicator 任务指示器，使用 Fail 执行任务失败的函数，使用 Success 执行任务成功的函数，并将失败方法置nil
type TaskIndicator struct {
	successDo func()
	failedDo  func()
}

func NewTaskIndicator(success, fail func()) *TaskIndicator {
	return &TaskIndicator{
		successDo: success,
		failedDo:  fail,
	}
}

func (t *TaskIndicator) Fail() {
	if t.failedDo != nil {
		t.failedDo()
	}
}

func (t *TaskIndicator) Success() {
	if t.successDo != nil {
		t.successDo()
	}
	t.failedDo = nil
}

// IsTerminalInput 是否通过终端拿取输入
func IsTerminalInput() bool {
	_, err := unix.IoctlGetTermios(unix.Stdin, unix.TCGETS)
	return err == nil
}
