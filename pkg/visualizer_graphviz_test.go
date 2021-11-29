package pkg

import (
	"fmt"
	"strings"
	"testing"
)

func TestGraphvizOutput(t *testing.T) {
	machineUnderTest := NewMachine(
		[]TransitionDesc{
			{Name: "open", Sources: []string{"closed"}, Destination: "open"},
			{Name: "close", Sources: []string{"open"}, Destination: "closed"},
			{Name: "part-close", Sources: []string{"intermediate"}, Destination: "closed"},
		},
		map[string]Callback{},
	)

	i := machineUnderTest.NewInstance("closed")

	got := Visualize(machineUnderTest, i)

	wanted := `
digraph fsm {
    "closed" -> "open" [ label = "open" ];
    "intermediate" -> "closed" [ label = "part-close" ];
    "open" -> "closed" [ label = "close" ];

    "closed";
    "intermediate";
    "open";
}`
	normalizedGot := strings.ReplaceAll(got, "\n", "")
	normalizedWanted := strings.ReplaceAll(wanted, "\n", "")
	if normalizedGot != normalizedWanted {
		t.Errorf("build graphivz graph failed. \nwanted \n%s\nand got \n%s\n", wanted, got)
		fmt.Println([]byte(normalizedGot))
		fmt.Println([]byte(normalizedWanted))
	}
}
