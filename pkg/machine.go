package pkg

import "strings"

// Machine is the state machine descriptor that holds the blueprint of the FSM.
//
// It has to be created with NewMachine to function properly.
type Machine struct {
	// transitions maps source states via a transition to destination states.
	transitions map[transitionKey]string

	// callbacks maps events and targets to callback functions.
	callbacks map[callbackKey]Callback
}

func NewMachine(transitions []TransitionDesc, callbacks map[string]Callback) *Machine {
	machine := &Machine{
		transitions: make(map[transitionKey]string),
		callbacks:   make(map[callbackKey]Callback),
	}

	// Build transition map and store sets of all events and states.
	allTransitions := make(map[string]bool)
	allStates := make(map[string]bool)
	for _, transition := range transitions {
		for _, source := range transition.Sources {
			transitionKey := transitionKey{transition.Name, source}
			machine.transitions[transitionKey] = transition.Destination
			allStates[source] = true
			allStates[transition.Destination] = true
		}
		allTransitions[transition.Name] = true
	}

	// Map all callbacks to transitions/states.
	for name, callback := range callbacks {
		callbackType := callbackNone
		var target string

		switch {
		case strings.HasPrefix(name, "before_"):
			target = strings.TrimPrefix(name, "before_")
			if target == "transition" {
				target = ""
				callbackType = callbackBeforeTransition
			} else if _, ok := allTransitions[target]; ok {
				callbackType = callbackBeforeTransition
			}
		case strings.HasPrefix(name, "leave_"):
			target = strings.TrimPrefix(name, "leave_")
			if target == "state" {
				target = ""
				callbackType = callbackLeaveState
			} else if _, ok := allStates[target]; ok {
				callbackType = callbackLeaveState
			}
		case strings.HasPrefix(name, "enter_"):
			target = strings.TrimPrefix(name, "enter_")
			if target == "state" {
				target = ""
				callbackType = callbackEnterState
			} else if _, ok := allStates[target]; ok {
				callbackType = callbackEnterState
			}
		case strings.HasPrefix(name, "after_"):
			target = strings.TrimPrefix(name, "after_")
			if target == "transition" {
				target = ""
				callbackType = callbackAfterTransition
			} else if _, ok := allTransitions[target]; ok {
				callbackType = callbackAfterTransition
			}
		default:
			target = name
			if _, ok := allStates[target]; ok {
				callbackType = callbackEnterState
			} else if _, ok := allTransitions[target]; ok {
				callbackType = callbackAfterTransition
			}
		}

		if callbackType != callbackNone {
			machine.callbacks[callbackKey{target, callbackType}] = callback
		}
	}

	return machine
}

func (machine *Machine) NewInstance(initial string) *Instance {
	return &Instance{
		current:         initial,
		transitionerObj: &transitionerStruct{},
		metadata:        make(map[string]interface{}),
	}
}
