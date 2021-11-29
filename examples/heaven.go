//go:build ignore

package main

import (
	"fmt"

	fsm "github.com/snapp-incubator/fsm/pkg"
)

func main() {
	machine := fsm.NewMachine(
		[]fsm.TransitionDesc{
			{Name: "open", Sources: []string{"closed"}, Destination: "open"},
			{Name: "close", Sources: []string{"open"}, Destination: "closed"},
		},
		map[string]fsm.Callback{
			"enter_state": func(e *fsm.Transition) { fmt.Printf("The door to heaven is %s\n", e.Dst) },
		},
	)

	heaven := machine.NewInstance("closed")

	err := heaven.Transition(machine, "open")
	if err != nil {
		fmt.Println(err)
	}

	err = heaven.Transition(machine, "close")
	if err != nil {
		fmt.Println(err)
	}
}
