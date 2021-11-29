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
		map[string]fsm.Callback{},
	)

	instance1 := machine.NewInstance("closed")
	instance2 := machine.NewInstance("open")

	fmt.Printf("instance1 state: %s\n", instance1.Current())
	fmt.Printf("instance2 state: %s\n", instance2.Current())

	err := instance1.Transition(machine, "open")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("instance1 state: %s\n", instance1.Current())
	fmt.Printf("instance2 state: %s\n", instance2.Current())

	err = instance1.Transition(machine, "close")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("instance1 state: %s\n", instance1.Current())
	fmt.Printf("instance2 state: %s\n", instance2.Current())
}
