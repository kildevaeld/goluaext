package tasks

import (
	"github.com/aarzilli/golua/lua"
	"github.com/kildevaeld/goluaext/runloop"
	"github.com/robertkrimen/otto"
	"github.com/stevedonovan/luar"
)

type simple_task struct {
	id     int64
	err    error
	result *luar.LuaObject
	call   *luar.LuaObject
}

func (self *simple_task) SetID(id int64) { self.id = id }
func (self *simple_task) GetID() int64   { return self.id }

func (self *simple_task) Execute(vm *lua.State, loop *runloop.Loop) error {

	var arguments []interface{}

	if self.err != nil {
		/*e, err := vm.Call(`new Error`, nil, self.err.Error())
		if err != nil {
			return err
		}*/
		panic(self.err.Error)

		arguments = append(arguments, e)
	} else {
		arguments = append(arguments, otto.NullValue())
	}

	arguments = append(arguments, self.result)

	//arguments = append([]interface{}{self.call}, arguments...)
	/*if _, err := self.call.Call(otto.NullValue(), arguments...); err != nil {
		return err
	}*/
	/*if _, err := vm.Call(`Function.call.call`, nil, arguments...); err != nil {
		return err
	}**/

	return nil
}

func (self *simple_task) Cancel() {

}

func SimpleTask(loop *runloop.Loop, cb *luar.LuaObject, worker func(task runloop.Task) (interface{}, error)) {

	task := &simple_task{call: cb}
	loop.Add(task)

	go func() {
		defer loop.Ready(task)
		var (
			i interface{}
			e error
			//v otto.Value
		)
		i, e = worker(task)
		task.err = e
		if i == nil {

			//task.result = otto.UndefinedValue()
		} else {
			task.result = luar.NewLuaObjectFromValue(loop.VM(), i)

			/*task.result, e = vm.ToValue(i)
			if e != nil {
				task.err = e
			}*/
		}

	}()

}
