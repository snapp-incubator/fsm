package pkg

type Transition struct {
	// Instance is an reference to the current FSM.
	Instance *Instance

	// Name is the transition name.
	Name string

	// Src is the state before the transition.
	Src string

	// Dst is the state after the transition.
	Dst string

	// Err is an optional error that can be returned from a callback.
	Err error

	// Args is an optional list of arguments passed to the callback.
	Args []interface{}

	// canceled is an internal flag set if the transition is canceled.
	canceled bool

	// async is an internal flag set if the transition should be asynchronous
	async bool
}

// Cancel can be called in before_<Transition> or leave_<STATE> to cancel the
// current transition before it happens. It takes an optional error, which will
// overwrite e.Err if set before.
func (t *Transition) Cancel(err ...error) {
	t.canceled = true

	if len(err) > 0 {
		t.Err = err[0]
	}
}

// Async can be called in leave_<STATE> to do an asynchronous state transition.
//
// The current state transition will be on hold in the old state until a final
// call to Transition is made. This will complete the transition and possibly
// call the other callbacks.
func (t *Transition) Async() {
	t.async = true
}

// TransitionDesc represents an event when initializing the FSM.
//
// The event can have one or more source states that is valid for performing
// the transition. If the FSM is in one of the source states it will end up in
// the specified destination state, calling all defined callbacks as it goes.
type TransitionDesc struct {
	// Name is the event name used when calling for a transition.
	Name string

	// Sources is a slice of source states that the FSM must be in to perform a state transition.
	Sources []string

	// Destination is the destination state that the FSM will be in if the transition succeeds.
	Destination string
}

// transitionKey is a struct key used for storing the transition map.
type transitionKey struct {
	// name of the transition that the keys refers to.
	name string

	// source from where the transition can transition.
	source string
}

// transitioner is an interface for the FSM's transition function.
type transitioner interface {
	transition(*Instance) error
}

// transitionerStruct is the default implementation of the transitioner
// interface. Other implementations can be swapped in for testing.
type transitionerStruct struct{}

// Transition completes an asynchronous state change.
//
// The callback for leave_<STATE> must previously have called Async on its
// event to have initiated an asynchronous state transition.
func (t transitionerStruct) transition(instance *Instance) error {
	if instance.transition == nil {
		return NotInTransitionError{}
	}
	instance.transition()
	instance.transition = nil
	return nil
}
