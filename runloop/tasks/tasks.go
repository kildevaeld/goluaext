package tasks

import (
	"errors"

	"github.com/aarzilli/golua/lua"
	"github.com/kildevaeld/goluaext/runloop"
	"github.com/robertkrimen/otto"
	"github.com/stevedonovan/luar"
)

// IdleTask is designed to sit in a loop and keep it active, without doing any
// work.
type IdleTask struct {
	ID int64
}

// NewIdleTask creates a new IdleTask object.
func NewIdleTask() *IdleTask {
	return &IdleTask{}
}

// SetID sets the ID of an IdleTask.
func (i *IdleTask) SetID(ID int64) { i.ID = ID }

// GetID gets the ID of an IdleTask.
func (i IdleTask) GetID() int64 { return i.ID }

// Cancel does nothing on an IdleTask, as there's nothing to clean up.
func (i IdleTask) Cancel() {}

// Execute always returns an error for an IdleTask, as it should never
// actually be run.
func (i IdleTask) Execute(vm *lua.State, l *runloop.Loop) error {
	return errors.New("Idle task should never execute")
}

// EvalTask schedules running an otto.Script. It has two channels for
// communicating the result of the execution.
type EvalTask struct {
	ID        int64
	Script    interface{}
	Value     chan *luar.LuaObject
	Error     chan error
	SoftError bool
}

// NewEvalTask creates a new EvalTask for a given otto.Script, creating two
// buffered channels for the response.
func NewEvalTask(s interface{}) *EvalTask {
	return &EvalTask{
		Script: s,
		Value:  make(chan *luar.LuaObject, 1),
		Error:  make(chan error, 1),
	}
}

// SetID sets the ID of an EvalTask.
func (e *EvalTask) SetID(ID int64) { e.ID = ID }

// GetID gets the ID of an EvalTask.
func (e EvalTask) GetID() int64 { return e.ID }

// Cancel does nothing for an EvalTask, as there's nothing to clean up.
func (e EvalTask) Cancel() {}

// Execute runs the EvalTask's otto.Script in the vm provided, pushing the
// resultant return value and error (or nil) into the associated channels.
// If the execution results in an error, it will return that error.
func (e EvalTask) Execute(vm *lua.State, l *runloop.Loop) error {
	v, err := vm.DoString(e.Script.(string))
	e.Value <- v
	e.Error <- err

	if e.SoftError {
		return nil
	}

	return err
}

// CallTask schedules an *luar.LuaObject (which should be a function) to be called
// with a specific set of arguments. It has two channels for communicating the
// result of the call.
type CallTask struct {
	ID        int64
	Function  *luar.LuaObject
	Args      []interface{}
	Value     chan *luar.LuaObject
	Error     chan error
	SoftError bool
}

// NewCallTask creates a new CallTask object for a given *luar.LuaObject (which
// should be a function) and set of arguments, creating two buffered channels
// for the response.
func NewCallTask(fn *luar.LuaObject, args ...interface{}) *CallTask {
	return &CallTask{
		Function: fn,
		Args:     args,
		Value:    make(chan *luar.LuaObject, 1),
		Error:    make(chan error, 1),
	}
}

// SetID sets the ID of a CallTask.
func (c *CallTask) SetID(ID int64) { c.ID = ID }

// GetID gets the ID of a CallTask.
func (c CallTask) GetID() int64 { return c.ID }

// Cancel does nothing for a CallTask, as there's nothing to clean up.
func (c CallTask) Cancel() {}

// Execute calls the associated function (not necessarily in the given vm),
// pushing the resultant return value and error (or nil) into the associated
// channels. If the call results in an error, it will return that error.
func (c CallTask) Execute(vm *lua.State, l *runloop.Loop) error {
	v, err := c.Function.Call(otto.NullValue(), c.Args...)
	c.Value <- v
	c.Error <- err

	if c.SoftError {
		return nil
	}

	return err
}
