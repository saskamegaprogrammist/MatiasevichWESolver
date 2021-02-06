package solver

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

const (
	FILENAME = "eq_graph_"
	GraphEXT = ".dot"
	PicEXT   = ".png"
)

type Writer struct {
	writer   *bufio.Writer
	file     *os.File
	filename string
}

func createFileName() string {
	return fmt.Sprintf("%s%d", FILENAME, time.Now().UnixNano())
}

func (writer *Writer) GetGraphFilename() string {
	return fmt.Sprintf("%s%s", writer.filename, GraphEXT)
}

func (writer *Writer) GetPicFilename() string {
	return fmt.Sprintf("%s%s", writer.filename, PicEXT)
}

func (writer *Writer) Init() error {
	var err error
	var file *os.File
	writer.filename = createFileName()
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
