package pkg

import (
	"errors"
	"sync"
)

type Instance struct {
	// current is the state that the FSM is currently in.
	current string

	// transition is the internal transition functions used either directly
	// or when Transition is called in an asynchronous state transition.
	transition func()
	// transitionerObj calls the FSM's transition() function.
	transitionerObj transitioner

	// stateMu guards access to the current state.
	stateMu sync.RWMutex
	// eventMu guards access to Event() and Transition().
	eventMu sync.Mutex
	// metadata can be used to store and load data that maybe used across events
	// use methods SetMetadata() and Metadata() to store and load data
	metadata map[string]interface{}

	metadataMu sync.RWMutex
}

// Current returns the current state of the FSM.
func (f *Instance) Current() string {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()

	return f.current
}

// Is returns true if state is the current state.
func (f *Instance) Is(state string) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()

	return state == f.current
}

// SetState allows the user to move to the given state from current state.
// The call does not trigger any callbacks, if defined.
func (f *Instance) SetState(state string) {
	f.stateMu.Lock()
	defer f.stateMu.Unlock()

	f.current = state
}

// Can returns true if event can occur in the current state.
func (f *Instance) Can(machine *Machine, event string) bool {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()

	_, ok := machine.transitions[transitionKey{event, f.current}]

	return ok && (f.transition == nil)
}

// AvailableTransitions returns a list of transitions available in the current state.
func (f *Instance) AvailableTransitions(machine *Machine) []string {
	f.stateMu.RLock()
	defer f.stateMu.RUnlock()

	var transitions []string
	for key := range machine.transitions {
		if key.source == f.current {
			transitions = append(transitions, key.name)
		}
	}

	return transitions
}

// SetMetadata stores the dataValue in metadata indexing it with key.
func (f *Instance) SetMetadata(key string, dataValue interface{}) {
	f.metadataMu.Lock()
	defer f.metadataMu.Unlock()
	f.metadata[key] = dataValue
}

// GetMetadata returns the value stored in metadata.
func (f *Instance) GetMetadata(key string) (interface{}, bool) {
	f.metadataMu.RLock()
	defer f.metadataMu.RUnlock()

	dataElement, ok := f.metadata[key]

	return dataElement, ok
}

// Transition initiates a state transition with the named event.
//
// The call takes a variable number of arguments that will be passed to the
// callback, if defined.
//
// It will return nil if the state change is ok or one of these errors:
//
// - event X inappropriate because previous transition did not complete
//
// - event X inappropriate in current state Y
//
// - event X does not exist
//
// - internal error on state transition
//
// The last error should never occur in this situation and is a sign of an
// internal bug.
func (f *Instance) Transition(machine *Machine, name string, args ...interface{}) error {
	f.eventMu.Lock()
	defer f.eventMu.Unlock()

	f.stateMu.RLock()
	defer f.stateMu.RUnlock()

	if f.transition != nil {
		return InTransitionError{name}
	}

	dst, ok := machine.transitions[transitionKey{name, f.current}]
	if !ok {
		for transitionkey := range machine.transitions {
			if transitionkey.name == name {
				return InvalidEventError{name, f.current}
			}
		}

		return UnknownEventError{name}
	}

	e := &Transition{f, name, f.current, dst, nil, args, false, false}

	err := f.beforeEventCallbacks(machine, e)
	if err != nil {
		return err
	}

	if f.current == dst {
		f.afterEventCallbacks(machine, e)

		return NoTransitionError{e.Err}
	}

	// Setup the transition, call it later.
	f.transition = func() {
		f.stateMu.Lock()
		f.current = dst
		f.stateMu.Unlock()

		f.enterStateCallbacks(machine, e)
		f.afterEventCallbacks(machine, e)
	}

	if err = f.leaveStateCallbacks(machine, e); err != nil {
		if ok := errors.As(err, new(CanceledError)); ok {
			f.transition = nil
		}

		return err
	}

	// Perform the rest of the transition, if not asynchronous.
	f.stateMu.RUnlock()
	defer f.stateMu.RLock()

	if err := f.doTransition(); err != nil {
		return InternalError{}
	}

	return e.Err
}

// doTransition wraps transitioner.transition.
func (f *Instance) doTransition() error {
	return f.transitionerObj.transition(f)
}
