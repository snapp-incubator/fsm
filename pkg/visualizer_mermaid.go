package pkg

import (
	"bytes"
	"fmt"
)

const highlightingColor = "#00AA00"

// MermaidDiagramType the type of the mermaid diagram type
type MermaidDiagramType string

const (
	// FlowChart the diagram type for output in flowchart style (https://mermaid-js.github.io/mermaid/#/flowchart) (including current state)
	FlowChart MermaidDiagramType = "flowChart"
	// StateDiagram the diagram type for output in stateDiagram style (https://mermaid-js.github.io/mermaid/#/stateDiagram)
	StateDiagram MermaidDiagramType = "stateDiagram"
)

// VisualizeForMermaidWithGraphType outputs a visualization of a FSM in Mermaid format as specified by the graphType.
func VisualizeForMermaidWithGraphType(machine *Machine, fsm *Instance, graphType MermaidDiagramType) (string, error) {
	switch graphType {
	case FlowChart:
		return visualizeForMermaidAsFlowChart(machine, fsm), nil
	case StateDiagram:
		return visualizeForMermaidAsStateDiagram(machine, fsm), nil
	default:
		return "", fmt.Errorf("unknown MermaidDiagramType: %s", graphType)
	}
}

func visualizeForMermaidAsStateDiagram(machine *Machine, fsm *Instance) string {
	var buf bytes.Buffer

	sortedTransitionKeys := getSortedTransitionKeys(machine.transitions)

	buf.WriteString("stateDiagram-v2\n")
	buf.WriteString(fmt.Sprintln(`    [*] -->`, fsm.current))

	for _, k := range sortedTransitionKeys {
		v := machine.transitions[k]
		buf.WriteString(fmt.Sprintf(`    %s --> %s: %s`, k.source, v, k.name))
		buf.WriteString("\n")
	}

	return buf.String()
}

// visualizeForMermaidAsFlowChart outputs a visualization of a FSM in Mermaid format (including highlighting of current state).
func visualizeForMermaidAsFlowChart(machine *Machine, fsm *Instance) string {
	var buf bytes.Buffer

	sortedTransitionKeys := getSortedTransitionKeys(machine.transitions)
	sortedStates, statesToIDMap := getSortedStates(machine.transitions)

	writeFlowChartGraphType(&buf)
	writeFlowChartStates(&buf, sortedStates, statesToIDMap)
	writeFlowChartTransitions(&buf, machine.transitions, sortedTransitionKeys, statesToIDMap)
	writeFlowChartHighlightCurrent(&buf, fsm.current, statesToIDMap)

	return buf.String()
}

func writeFlowChartGraphType(buf *bytes.Buffer) {
	buf.WriteString("graph LR\n")
}

func writeFlowChartStates(buf *bytes.Buffer, sortedStates []string, statesToIDMap map[string]string) {
	for _, state := range sortedStates {
		buf.WriteString(fmt.Sprintf(`    %s[%s]`, statesToIDMap[state], state))
		buf.WriteString("\n")
	}

	buf.WriteString("\n")
}

func writeFlowChartTransitions(buf *bytes.Buffer, transitions map[transitionKey]string, sortedTransitionKeys []transitionKey, statesToIDMap map[string]string) {
	for _, transition := range sortedTransitionKeys {
		target := transitions[transition]
		buf.WriteString(fmt.Sprintf(`    %s --> |%s| %s`, statesToIDMap[transition.source], transition.name, statesToIDMap[target]))
		buf.WriteString("\n")
	}
	buf.WriteString("\n")
}

func writeFlowChartHighlightCurrent(buf *bytes.Buffer, current string, statesToIDMap map[string]string) {
	buf.WriteString(fmt.Sprintf(`    style %s fill:%s`, statesToIDMap[current], highlightingColor))
	buf.WriteString("\n")
}
