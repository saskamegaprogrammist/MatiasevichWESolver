package solver

import (
	"fmt"
	"github.com/goccy/go-graphviz"
	"github.com/saskamegaprogrammist/MatiasevichWESolver/solver/equation/symbol"
	"io/ioutil"
)

type DotWriter struct {
	writer Writer
}

func (dotWriter *DotWriter) Init(mode string, eq string, outputDir string) error {
	err := dotWriter.writer.Init(mode, eq, outputDir)
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

func (dotWriter *DotWriter) EndDOTDescription(makePng bool) error {
	err := dotWriter.writer.Write("}")
	if err != nil {
		return fmt.Errorf("error ending DOT description: %v", err)
	}
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

func getEdgeLabel(symbol *symbol.Symbol, newSymbols []symbol.Symbol) string {
	label := fmt.Sprintf("%s->", (*symbol).Value())
	for _, sym := range newSymbols {
		label += sym.Value()
	}
	return label
}

func (dotWriter *DotWriter) WriteEdge(from *Node, to *Node) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s;\n", from.number, to.number))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteLabelEdge(from *Node, to *Node, symbol *symbol.Symbol, newSymbols []symbol.Symbol) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s[label=\"%s\"];\n", from.number, to.number, getEdgeLabel(symbol, newSymbols)))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteLabelEdgeSystem(from *NodeSystem, to *NodeSystem, symbol *symbol.Symbol, newSymbols []symbol.Symbol) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s[label=\"%s\"];\n", from.number, to.number, getEdgeLabel(symbol, newSymbols)))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteLabelEdgeBoldSystem(from *NodeSystem, to *NodeSystem) error {
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

func (dotWriter *DotWriter) WriteInfoEdgeWithLabel(from *Node, to InfoNode, label string) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s[label=\"%s\"];\n", from.number, to.GetNumber(), label))
	if err != nil {
		return fmt.Errorf("error describing edge: %v", err)
	}
	return nil
}

func (dotWriter *DotWriter) WriteInfoEdgeSystem(from *NodeSystem, to InfoNode) error {
	err := dotWriter.writer.Write(fmt.Sprintf("     %s -> %s;\n", from.number, to.GetNumber()))
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

func (dotWriter *DotWriter) WriteDottedEdgeSystem(from *NodeSystem, to *NodeSystem) error {
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

func (dotWriter *DotWriter) WriteNodeSystem(nodeSystem *NodeSystem) error {
	err := dotWriter.writer.Write(fmt.Sprintf("    %s [label=\"%s\"];\n", nodeSystem.number, nodeSystem.Value.String()))
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
