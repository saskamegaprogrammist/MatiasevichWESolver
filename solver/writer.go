package solver

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const (
	FILENAME = "eq_graph_"
	GraphEXT = ".dot"
	PicEXT   = ".png"
	DOT      = "."
)

var number = 0

type Writer struct {
	writer    *bufio.Writer
	file      *os.File
	filename  string
	outputDir string
}

func (writer *Writer) createFileName(name string, defaultFilename bool) {
	var newName string
	if defaultFilename {
		parsed := writer.parseName(name)
		if parsed == "" {
			newName = randStr(10)
		} else {
			newName = parsed
		}
	} else {
		newName = FILENAME + name
	}
	writer.filename = writer.outputDir + string(os.PathSeparator) + newName
}

func (writer *Writer) parseName(defaultName string) string {
	splitted := strings.Split(defaultName, string(os.PathSeparator))
	newSplitted := strings.Split(splitted[len(splitted)-1], DOT)
	return newSplitted[0]
}

func (writer *Writer) modifyFileName() {
	if number > 0 {
		ind := writer.filename[len(writer.filename)-1:]
		indI, _ := strconv.Atoi(ind)
		writer.filename = fmt.Sprintf("%s%d", writer.filename[:len(writer.filename)-1], indI+1)
	} else {
		writer.filename = fmt.Sprintf("%s_%d", writer.filename, number)
	}
}

func (writer *Writer) GetGraphFilename() string {
	return fmt.Sprintf("%s%s", writer.filename, GraphEXT)
}

func (writer *Writer) GetPicFilename() string {
	return fmt.Sprintf("%s%s", writer.filename, PicEXT)
}

func (writer *Writer) Init(outName string, defaultFilename bool, outputDir string) error {
	var err, fErr error
	var file *os.File
	writer.outputDir = outputDir
	writer.createFileName(outName, defaultFilename)
	_, fErr = os.Stat(writer.GetGraphFilename())
	for fErr == nil {
		writer.modifyFileName()
		_, fErr = os.Stat(writer.GetGraphFilename())
		number++
	}
	number = 1
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
