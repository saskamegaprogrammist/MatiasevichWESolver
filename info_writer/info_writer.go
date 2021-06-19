package info_writer

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	INFONAME = "info.txt"
)

type InfoWriter struct {
	writer *bufio.Writer
	file   *os.File
}

func (infoWriter *InfoWriter) Init() error {
	var err error
	var file *os.File
	file, err = os.Create(INFONAME)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	infoWriter.file = file
	infoWriter.writer = bufio.NewWriter(infoWriter.file)
	return nil
}

func (infoWriter *InfoWriter) WriteEquation(equation string) error {
	_, err := infoWriter.writer.WriteString(fmt.Sprintf("%s\n", equation))
	if err != nil {
		return fmt.Errorf("error wriring to info writer: %v", err)
	}
	return nil
}

func (infoWriter *InfoWriter) WriteFormat(mode string) error {
	_, err := infoWriter.writer.WriteString(fmt.Sprintf("mode: %s\n", mode))
	if err != nil {
		return fmt.Errorf("error wriring to info writer: %v", err)
	}
	return nil
}

func (infoWriter *InfoWriter) WriteTime(time time.Duration) error {
	_, err := infoWriter.writer.WriteString(fmt.Sprintf("took time: %v\n", time))
	if err != nil {
		return fmt.Errorf("error wriring to info writer: %v", err)
	}
	return nil
}

func (infoWriter *InfoWriter) WriteSolution(solution string) error {
	_, err := infoWriter.writer.WriteString(fmt.Sprintf("solution: %s\n\n", solution))
	if err != nil {
		return fmt.Errorf("error wriring to info writer: %v", err)
	}
	return nil
}

func (infoWriter *InfoWriter) Write(str string) error {
	_, err := infoWriter.writer.WriteString(str)
	if err != nil {
		return fmt.Errorf("error wriring to info writer: %v", err)
	}
	return nil
}

func (infoWriter *InfoWriter) Flush() error {
	err := infoWriter.writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing to info writer: %v", err)
	}
	return nil
}
