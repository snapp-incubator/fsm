package pkg

// Callback is a function type that callbacks should use.
// Transition is the current transition as the callback happens.
type Callback func(*Transition)

type callbackType uint8

const (
	callbackNone callbackType = iota
	callbackBeforeTransition
	callbackLeaveState
	callbackEnterState
	callbackAfterTransition
)

// cKey is a struct key used for keeping the callbacks mapped to a target.
type callbackKey struct {
	// target is either the name of a state or an event depending on which
	// callback type the key refers to. It can also be "" for a non-targeted
	// callback like before_event.
	target string

	// callbackType is the situation when the callback will be run.
	callbackType callbackType
}

// beforeEventCallbacks calls the before_ callbacks, first the named then the general version.
func (f *Instance) beforeEventCallbacks(machine *Machine, t *Transition) error {
	if fn, ok := machine.callbacks[callbackKey{t.Name, callbackBeforeTransition}]; ok {
		fn(t)
		if t.canceled {
			return CanceledError{t.Err}
		}
	}
	if fn, ok := machine.callbacks[callbackKey{"", callbackBeforeTransition}]; ok {
		fn(t)
		if t.canceled {
			return CanceledError{t.Err}
		}
	}
	return nil
}

// leaveStateCallbacks calls the leave_ callbacks, first the named then the general version.
func (f *Instance) leaveStateCallbacks(machine *Machine, e *Transition) error {
	if fn, ok := machine.callbacks[callbackKey{f.current, callbackLeaveState}]; ok {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		} else if e.async {
			return AsyncError{e.Err}
		}
	}
	if fn, ok := machine.callbacks[callbackKey{"", callbackLeaveState}]; ok {
		fn(e)
		if e.canceled {
			return CanceledError{e.Err}
		} else if e.async {
			return AsyncError{e.Err}
		}
	}
	return nil
}

// enterStateCallbacks calls the enter_ callbacks, first the named then the general version.
func (f *Instance) enterStateCallbacks(machine *Machine, e *Transition) {
	if fn, ok := machine.callbacks[callbackKey{f.current, callbackEnterState}]; ok {
		fn(e)
	}
	if fn, ok := machine.callbacks[callbackKey{"", callbackEnterState}]; ok {
		fn(e)
	}
}

// afterEventCallbacks calls the after_ callbacks, first the named then the general version.
func (f *Instance) afterEventCallbacks(machine *Machine, e *Transition) {
	if fn, ok := machine.callbacks[callbackKey{e.Name, callbackAfterTransition}]; ok {
		fn(e)
	}
	if fn, ok := machine.callbacks[callbackKey{"", callbackAfterTransition}]; ok {
		fn(e)
	}
}
