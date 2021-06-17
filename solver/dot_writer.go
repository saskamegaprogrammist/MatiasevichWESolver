package solver

import (
	"fmt"
	"github.com/goccy/go-graphviz"
	"io/ioutil"
)

type DotWriter struct {
	writer Writer
}

func (dotWriter *DotWriter) Init(outName string, defaultFilename bool, outputDir string) error {
	err := dotWriter.writer.Init(outName, defaultFilename, outputDir)
	if err != nil {
		return fmt.Errorf("error initing writer: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) StartDOTDescription() error {
	err := dotWriter.writer.Write("digraph word_eq {\n")
	if err != nil {
		return fmt.Errorf("error starting DOT description: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) EndDOTDescription() error {
	err := dotWriter.writer.Write("}")
	if err != nil {
		return fmt.Errorf("error ending DOT description: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) CreateFiles(makePng bool) error {
	var err error
	err = dotWriter.writer.Flush()
	if err != nil {
		return fmt.Errorf("error flushing DOT description: %v", err)
	}
	if makePng {
		err = dotWriter.CreatePNG()
		if err != nil {
			return fmt.Errorf("error creating png: %v", err)
		}
	}
	return nil
}

func (dotWriter *DotWriter) WriteEdge(from *Node, to *Node) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s;\n", from.number, to.number))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteLabelEdge(from *Node, to *Node, label string) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s[label=\"%s\"];\n", from.number, to.number, label))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteLabelEdgeBold(from *Node, to *Node) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s[style=bold][label=splitting];\n", from.number, to.number))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteInfoEdge(from *Node, to InfoNode) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s;\n", from.number, to.GetNumber()))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteInfoEdgeWithLabel(from *Node, to *FalseNode) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s[label=\"%s\"];\n", from.number, to.GetNumber(), to.GetInfoLabel()))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteDottedEdge(from *Node, to *Node) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s [style=dotted];\n", from.number, to.number))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteNode(node *Node) error {
	err := dotWriter.writer.Write(fmt.Sprintf("    %s [label=\"%s\"];\n", node.number, node.value.String()))
	if err != nil {
		return fmt.Errorf("error describing node: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteInfoNode(node InfoNode) error {
	err := dotWriter.writer.Write(fmt.Sprintf("    %s [label=\"%s\"];\n", node.GetNumber(), node.GetValue()))
	if err != nil {
		return fmt.Errorf("error describing node: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) CreatePNG() error {
	bytes, err := ioutil.ReadFile(dotWriter.writer.GetGraphFilename())
	if err != nil {
		return fmt.Errorf("error reading dot file: %v", err)
	}
	graph, err := graphviz.ParseBytes(bytes)
	if err != nil {
		return fmt.Errorf("error parsing dot file: %v", err)
	}
	if err := graphviz.New().RenderFilename(graph, graphviz.PNG, dotWriter.writer.GetPicFilename()); err != nil {
		return fmt.Errorf("error writing to png file: %v", err)
	}
	return nil
}
