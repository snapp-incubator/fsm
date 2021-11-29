//go:build ignore

package main

import (
	"fmt"

	fsm "github.com/snapp-incubator/fsm/pkg"
)

func main() {
	machine := fsm.NewMachine(
		[]fsm.TransitionDesc{
			{Name: "scan", Sources: []string{"idle"}, Destination: "scanning"},
			{Name: "working", Sources: []string{"scanning"}, Destination: "scanning"},
			{Name: "situation", Sources: []string{"scanning"}, Destination: "scanning"},
			{Name: "situation", Sources: []string{"idle"}, Destination: "idle"},
			{Name: "finish", Sources: []string{"scanning"}, Destination: "idle"},
		},
		map[string]fsm.Callback{
			"scan": func(t *fsm.Transition) {
				fmt.Println("after_scan: " + t.Instance.Current())
			},
			"working": func(t *fsm.Transition) {
				fmt.Println("working: " + t.Instance.Current())
			},
			"situation": func(t *fsm.Transition) {
				fmt.Println("situation: " + t.Instance.Current())
			},
			"finish": func(t *fsm.Transition) {
				fmt.Println("finish: " + t.Instance.Current())
			},
		},
	)

	instance := machine.NewInstance("idle")

	fmt.Println(instance.Current())

	err := instance.Transition(machine, "scan")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("1:" + instance.Current())

	err = instance.Transition(machine, "working")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("2:" + instance.Current())

	err = instance.Transition(machine, "situation")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("3:" + instance.Current())

	err = instance.Transition(machine, "finish")
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("4:" + instance.Current())

}
