package auth

import "github.com/oracle/oci-go-sdk/common"

//dispatcherModifier gives ability to modify a HTTPRequestDispatcher before use.
type dispatcherModifier struct {
	modifiers []func(common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error)
}

//newDispatcherModifier creates a new dispatcherModifier with optional initial modifier (may be nil).
func newDispatcherModifier(modifier func(common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error)) *dispatcherModifier {
	dispatcherModifier := &dispatcherModifier{
		modifiers: make([]func(common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error), 0),
	}
	if modifier != nil {
		dispatcherModifier.QueueModifier(modifier)
	}
	return dispatcherModifier
}

//QueueModifier queues up a new modifier
func (c *dispatcherModifier) QueueModifier(modifier func(common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error)) {
	c.modifiers = append(c.modifiers, modifier)
}

//Modify the provided HTTPRequestDispatcher with this modifier, and return the result, or error if something goes wrong
func (c *dispatcherModifier) Modify(dispatcher common.HTTPRequestDispatcher) (common.HTTPRequestDispatcher, error) {
	if len(c.modifiers) > 0 {
		for _, modifier := range c.modifiers {
			var err error
			if dispatcher, err = modifier(dispatcher); err != nil {
				common.Debugf("An error occurred when attempting to modify the dispatcher. Error was: %s", err.Error())
				return nil, err
			}
		}
	}
	return dispatcher, nil
}
