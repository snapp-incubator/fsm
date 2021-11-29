//go:build ignore

package main

import (
	"fmt"

	fsm "github.com/snapp-incubator/fsm/pkg"
)

func main() {
	counter := 0

	machine := fsm.NewMachine(
		[]fsm.TransitionDesc{
			{Name: "publish", Sources: []string{"idle"}, Destination: "idle"},
			{Name: "subscribe", Sources: []string{"idle"}, Destination: "idle"},
		},
		map[string]fsm.Callback{
			"publish": func(e *fsm.Transition) {
				msg := fmt.Sprintf("counter:%d", counter)
				e.Instance.SetMetadata("message", msg)
				fmt.Println("published data")
				counter++
			},
			"subscribe": func(e *fsm.Transition) {
				message, ok := e.Instance.GetMetadata("message")
				if ok {
					fmt.Println("message = " + message.(string))
				}
			},
		},
	)

	instance1 := machine.NewInstance("idle")
	instance2 := machine.NewInstance("idle")

	fmt.Printf("instance1 state: %s\n", instance1.Current())
	fmt.Printf("instance2 state: %s\n", instance2.Current())

	instance1.Transition(machine, "publish")
	instance2.Transition(machine, "publish")

	fmt.Printf("instance1 state: %s\n", instance1.Current())
	fmt.Printf("instance2 state: %s\n", instance2.Current())

	instance1.Transition(machine, "subscribe")
	instance2.Transition(machine, "subscribe")

	fmt.Printf("instance1 state: %s\n", instance1.Current())
	fmt.Printf("instance2 state: %s\n", instance2.Current())
}
