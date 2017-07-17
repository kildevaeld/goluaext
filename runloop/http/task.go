package http

import (
	"github.com/aarzilli/golua/lua"
	"github.com/kildevaeld/goluaext/runloop"
	"github.com/stevedonovan/luar"
)

type http_task struct {
	id int64
	//cb      *luar.LuaObject
	p       *luar.LuaObject
	err     error
	content luar.Map
}

func (t *http_task) SetID(id int64) { t.id = id }
func (t *http_task) GetID() int64   { return t.id }

func (t *http_task) Execute(vm *lua.State, loop *runloop.Loop) error {

	var (
		resolve *luar.LuaObject
		reject  *luar.LuaObject
		err     error
	)

	defer loop.Remove(t)

	if resolve, err = t.p.GetObject("resolve"); err != nil {
		return err
	}
	if reject, err = t.p.GetObject("reject"); err != nil {
		return err
	}
	defer func() {
		reject.Close()
		resolve.Close()
		t.p.Close()
		t.p = nil
	}()

	if t.err != nil {
		return reject.Call(nil, t.p, t.err.Error())
	}
	return resolve.Call(nil, t.p, t.content)
}

func (self *http_task) Cancel() {

}
