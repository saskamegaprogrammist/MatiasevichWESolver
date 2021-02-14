package solver

import (
	"bufio"
	"fmt"
	"os"
)

const (
	FILENAME = "eq_graph_"
	GraphEXT = ".dot"
	PicEXT   = ".png"
)

type Writer struct {
	writer    *bufio.Writer
	file      *os.File
	filename  string
	outputDir string
}

func (writer *Writer) createFileName(mode string, eq string) string {
	return fmt.Sprintf("%s%c%s%s_%s", writer.outputDir, os.PathSeparator, FILENAME, mode, eq)
}

func (writer *Writer) GetGraphFilename() string {
	return fmt.Sprintf("%s%s", writer.filename, GraphEXT)
}

func (writer *Writer) GetPicFilename() string {
	return fmt.Sprintf("%s%s", writer.filename, PicEXT)
}

func (writer *Writer) Init(mode string, eq string, outputDir string) error {
	var err error
	var file *os.File
	writer.outputDir = outputDir
	writer.filename = writer.createFileName(mode, eq)
	file, err = os.Create(writer.GetGraphFilename())
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}
	writer.file = file
	writer.writer = bufio.NewWriter(writer.file)
	return nil
}

func (writer *Writer) Write(str string) error {
	_, err := writer.writer.WriteString(str)
	if err != nil {
		return fmt.Errorf("error wriring to writer: %v", err)
	}
	return nil
}

func (writer *Writer) Flush() error {
	err := writer.writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing to writer: %v", err)
	}
	return nil
}
