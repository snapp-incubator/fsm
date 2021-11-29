//go:build ignore

package main

import (
	"fmt"

	fsm "github.com/snapp-incubator/fsm/pkg"
)

func main() {
	machine := fsm.NewMachine(
		[]fsm.TransitionDesc{
			{Name: "publish", Sources: []string{"idle"}, Destination: "idle"},
			{Name: "subscribe", Sources: []string{"idle"}, Destination: "idle"},
		},
		map[string]fsm.Callback{
			"publish": func(e *fsm.Transition) {
				e.Instance.SetMetadata("message", "hii")
				fmt.Println("published data")
			},
			"subscribe": func(e *fsm.Transition) {
				message, ok := e.Instance.GetMetadata("message")
				if ok {
					fmt.Println("message = " + message.(string))
				}

			},
		},
	)

	instance := machine.NewInstance("idle")

	fmt.Println(instance.Current())

	err := instance.Transition(machine, "publish")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(instance.Current())

	err = instance.Transition(machine, "subscribe")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(instance.Current())

}
