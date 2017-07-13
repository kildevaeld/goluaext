package goluaext

import (
	"errors"

	"github.com/aarzilli/golua/lua"
	"github.com/kildevaeld/go-args"
	"github.com/stevedonovan/luar"
)

var ErrCannotConvert = errors.New("cannot convert")

type Converter interface {
	Push(a args.Argument, state *lua.State) error
	Pop(state *lua.State, index int) (args.Argument, error)
}

var _converters map[args.Type]Converter

func RegisterConverter(hook args.Type, fn Converter) error {
	if _, ok := _converters[hook]; ok {
		return errors.New("converter already exists")
	}
	_converters[hook] = fn

	return nil
}

func convert_push(a args.Argument, state *lua.State) error {
	if converter, ok := _converters[a.Type()]; ok {
		return converter.Push(a, state)
	}

	return nil
}

func convert_pop(state *lua.State, index int) (args.Argument, error) {

	for _, c := range _converters {
		arg, err := c.Pop(state, index)
		if err == ErrCannotConvert {
			continue
		} else if err != nil {
			return nil, err
		} else if arg != nil {
			return arg, nil
		}

	}

	return nil, ErrCannotConvert
}

func init() {
	_converters = make(map[args.Type]Converter)
}

type call_argument struct {
	v *luar.LuaObject
}

func wrap_call(arguments args.Argument) interface{} {
	return func(state *lua.State) int {
		defer arguments.Free()
		a, err := LuaToArgument(state, 1)
		if err != nil {
			panic(err)
		}
		call := arguments.Value().(args.Call)
		if a, err = call.Call(a); err != nil {
			panic(err)
		} else if a != nil {
			PushArgument(state, a)
			return 1
		}

		return 0
	}
}

func (a *call_argument) Call(arguments args.Argument) (args.Argument, error) {

	var out []interface{}

	defer arguments.Free()
	val := arguments.Value()
	if arguments.Type() == args.CallType {
		val = wrap_call(arguments)
	}

	if err := a.v.Call(&out, val); err != nil {
		return nil, err
	}

	var arg args.Argument
	var err error
	if len(out) > 0 {
		if arg, err = args.NewArgument(out[0]); err != nil {
			return nil, err
		}
	}

	return arg, nil

}

func (a *call_argument) Free() {
	if a.v != nil {
		a.v.Close()
		a.v = nil
	}
}

func LuaToArgument(state *lua.State, i int) (args.Argument, error) {
	var arguments []args.Argument
	for {
		if state.IsTable(i) || state.IsGoStruct(i) {

			if arg, err := convert_pop(state, i); err == nil {
				arguments = append(arguments, arg)
				continue
			} else if err != ErrCannotConvert {
				return nil, err
			}

			var out luar.Map
			if err := luar.LuaToGo(state, i, &out); err != nil {
				return nil, err
			}
			arguments = append(arguments, args.NewArgumentOrNil(out))

		} else if state.IsBoolean(i) {
			arguments = append(arguments, args.NewArgumentOrNil(state.ToBoolean(i)))
		} else if state.IsString(i) {
			arguments = append(arguments, args.NewArgumentOrNil(state.ToString(i)))
		} else if state.IsNumber(i) {
			arguments = append(arguments, args.NewArgumentOrNil(state.ToNumber(i)))
		} else if state.IsNone(i) {
			break
		} else if state.IsFunction(i) {
			arguments = append(arguments, args.NewArgumentOrNil(&call_argument{luar.NewLuaObject(state, i)}))
		}

		i++
	}

	return args.NewArgument(arguments)
}

func PushArgument(state *lua.State, arg args.Argument) error {
	if arg == nil {
		return nil
	}

	switch arg.Type() {
	case args.StringType, args.NumberType, args.BoolType:
		luar.GoToLua(state, arg.Value())
	case args.CallType:
		luar.GoToLua(state, wrap_call(arg))
	case args.ArgumentListType:
		return PushArgumentList(state, arg.Value().(args.ArgumentList))
	case args.ArgumentMapType:
		return PushArgumentMap(state, arg.Value().(args.ArgumentMap))

	default:
		if err := convert_push(arg, state); err != nil {
			if err != ErrCannotConvert {
				return err
			}
		}
		return errors.New("invalid type")
	}

	return nil
}

func PushArgumentList(state *lua.State, arg args.ArgumentList) error {
	top := state.GetTop()
	l := len(arg)
	state.CreateTable(l, 0)
	for i := 0; i < l; i++ {
		if err := PushArgument(state, arg[i]); err != nil {
			state.SetTop(top)
			return err
		}
		state.RawSeti(-2, i)
	}

	return nil
}

func PushArgumentMap(state *lua.State, arg args.ArgumentMap) error {
	top := state.GetTop()
	state.CreateTable(0, len(arg))

	for k, v := range arg {
		if err := PushArgument(state, v); err != nil {
			state.SetTop(top)
			return err
		}
		state.SetField(-2, k)
	}

	return nil
}
