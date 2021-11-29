package pkg

import (
	"fmt"
	"sort"
)

// VisualizeType the type of the visualization
type VisualizeType string

const (
	// GRAPHVIZ the type for graphviz output (http://www.webgraphviz.com/)
	GRAPHVIZ VisualizeType = "graphviz"
	// MERMAID the type for mermaid output (https://mermaid-js.github.io/mermaid/#/stateDiagram) in the stateDiagram form
	MERMAID VisualizeType = "mermaid"
	// MermaidStateDiagram the type for mermaid output (https://mermaid-js.github.io/mermaid/#/stateDiagram) in the stateDiagram form
	MermaidStateDiagram VisualizeType = "mermaid-state-diagram"
	// MermaidFlowChart the type for mermaid output (https://mermaid-js.github.io/mermaid/#/flowchart) in the flow chart form
	MermaidFlowChart VisualizeType = "mermaid-flow-chart"
)

// VisualizeWithType outputs a visualization of a FSM in the desired format.
// If the type is not given it defaults to GRAPHVIZ
func VisualizeWithType(machine *Machine, fsm *Instance, visualizeType VisualizeType) (string, error) {
	switch visualizeType {
	case GRAPHVIZ:
		return Visualize(machine, fsm), nil
	case MERMAID:
		return VisualizeForMermaidWithGraphType(machine, fsm, StateDiagram)
	case MermaidStateDiagram:
		return VisualizeForMermaidWithGraphType(machine, fsm, StateDiagram)
	case MermaidFlowChart:
		return VisualizeForMermaidWithGraphType(machine, fsm, FlowChart)
	default:
		return "", fmt.Errorf("unknown VisualizeType: %s", visualizeType)
	}
}

func getSortedTransitionKeys(transitions map[transitionKey]string) []transitionKey {
	// we sort the key alphabetically to have a reproducible graph output
	sortedTransitionKeys := make([]transitionKey, 0)

	for transition := range transitions {
		sortedTransitionKeys = append(sortedTransitionKeys, transition)
	}
	sort.Slice(sortedTransitionKeys, func(i, j int) bool {
		if sortedTransitionKeys[i].source == sortedTransitionKeys[j].source {
			return sortedTransitionKeys[i].name < sortedTransitionKeys[j].name
		}
		return sortedTransitionKeys[i].source < sortedTransitionKeys[j].source
	})

	return sortedTransitionKeys
}

func getSortedStates(transitions map[transitionKey]string) ([]string, map[string]string) {
	statesToIDMap := make(map[string]string)
	for transition, target := range transitions {
		if _, ok := statesToIDMap[transition.source]; !ok {
			statesToIDMap[transition.source] = ""
		}
		if _, ok := statesToIDMap[target]; !ok {
			statesToIDMap[target] = ""
		}
	}

	sortedStates := make([]string, 0, len(statesToIDMap))
	for state := range statesToIDMap {
		sortedStates = append(sortedStates, state)
	}
	sort.Strings(sortedStates)

	for i, state := range sortedStates {
		statesToIDMap[state] = fmt.Sprintf("id%d", i)
	}
	return sortedStates, statesToIDMap
}
