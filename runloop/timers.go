package runloop

import (
	"time"

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

var minDelay = map[bool]int64{
	true:  10,
	false: 4,
}

type timerTask struct {
	id       int64
	timer    *time.Timer
	duration time.Duration
	interval bool
	call     *luar.LuaObject
	stopped  bool
}

func RegisterTimers(loop *Loop) error {
	state := loop.VM()

	newTimer := func(interval bool) func(state *lua.State) int {
		return func(state *lua.State) int {

			if !state.IsFunction(1) {
				panic("not function")
			}

			delay := int64(state.ToInteger(2))
			if delay < minDelay[interval] {
				delay = minDelay[interval]
			}

			t := &timerTask{
				duration: time.Duration(delay) * time.Millisecond,
				call:     luar.NewLuaObject(state, 1),
				interval: interval,
			}
			loop.Add(t)

			t.timer = time.AfterFunc(t.duration, func() {
				loop.Ready(t)
			})

			state.PushGoStruct(t)
			/*value, err := call.Otto.ToValue(t)
			if err != nil {
				panic(err)
			}*/

			return 1
		}
	}

	clearTimeout := func(call *lua.State) int {
		if !call.IsGoStruct(1) {
			panic("unreal")
		}
		if t, ok := call.ToGoStruct(1).(*timerTask); ok {
			t.stopped = true
			t.timer.Stop()
			loop.Remove(t)
		} else {
			panic("unreal 2")
		}

		return 0
	}

	luar.Register(state, "", luar.Map{
		"setTimeout":    newTimer(false),
		"clearTimeout":  clearTimeout,
		"setInterval":   newTimer(true),
		"clearInterval": clearTimeout,
	})

	return nil

}

func (t *timerTask) SetID(id int64) { t.id = id }
func (t *timerTask) GetID() int64   { return t.id }

func (t *timerTask) Execute(vm *lua.State, l *Loop) error {
	/*var arguments []interface{}

	if len(t.call.ArgumentList) > 2 {
		tmp := t.call.ArgumentList[2:]
		arguments = make([]interface{}, 2+len(tmp))

		for i, value := range tmp {
			arguments[i+2] = value
		}
	} else {
		arguments = make([]interface{}, 1)
	}*/

	/*arguments[0] = t.call.ArgumentList[0]

	if _, err := vm.Call(`Function.call.call`, nil, arguments...); err != nil {
		return err
	}*/

	/*if _, err := t.call.Argument(0).Call(otto.NullValue()); err != nil {
		return err
	}*/
	if err := t.call.Call(nil); err != nil {
		return err
	}
	if t.interval && !t.stopped {
		t.timer.Reset(t.duration)
		l.Add(t)
	}

	return nil
}

func (t *timerTask) Cancel() {
	t.timer.Stop()
}
