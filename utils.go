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
	l *lua.State
}

func wrap_call(arguments args.Argument) func(state *lua.State) int {
	return func(state *lua.State) int {
		defer arguments.Free()
		a, err := LuaToArgument(state, 1, false)
		if err != nil {
			panic(err)
		}
		call := arguments.Value().(args.Call)
		if a, err = call.Call(args.ArgumentList{a}); err != nil {
			panic(err)
		} else if a != nil {
			PushArgument(state, a)
			return 1
		}

		return 0
	}
}

func (a *call_argument) Call(arguments args.ArgumentList) (args.Argument, error) {

	top := a.l.GetTop()

	a.v.Push()
	defer a.l.SetTop(top)
	for _, arg := range arguments {
		if arg.Type() == args.CallType {
			a.l.PushGoFunction(wrap_call(arg))
		} else {
			if err := PushArgument(a.l, arg); err != nil {

				return nil, err
			}
		}
	}

	if err := a.l.Call(arguments.Len(), 1); err != nil {
		return nil, err
	}

	return LuaToArgument(a.l, 1, false)

}

func (a *call_argument) Free() {
	if a.v != nil {
		a.v.Close()
		a.v = nil
		a.l = nil
	}
}

func LuaToArgument(state *lua.State, i int, one bool) (args.Argument, error) {
	var arguments []args.Argument
	for {
		if state.IsTable(i) || state.IsGoStruct(i) {

			if arg, err := convert_pop(state, i); err == nil {
				arguments = append(arguments, arg)
				continue
			} else if err != ErrCannotConvert {
				return nil, err
			}
			// Is converted to ArgumentMap in Argument Constructor
			var out map[string]interface{}
			if err := luar.LuaToGo(state, i, &out); err != nil {
				return nil, err
			}
			arguments = append(arguments, args.Must(out))
		} else if state.IsBoolean(i) {
			arguments = append(arguments, args.Must(state.ToBoolean(i)))
		} else if state.IsString(i) {
			arguments = append(arguments, args.Must(state.ToString(i)))
		} else if state.IsNumber(i) {
			arguments = append(arguments, args.Must(state.ToNumber(i)))
		} else if state.IsNone(i) {
			//arguments = append(arguments, args.Undefined())
			break
		} else if state.IsFunction(i) {
			arguments = append(arguments, args.Must(&call_argument{luar.NewLuaObject(state, i), state}))
		} else {
			break
		}

		i++
		if one {
			return arguments[0], nil
		}
	}

	return args.New(arguments)
}

func isNumeric(a args.Argument) bool {
	switch a.Type() {
	case args.Int16Type, args.Int32Type, args.Int64Type, args.Int8Type, args.IntType,
		args.Uint16Type, args.Uint32Type, args.Uint64Type, args.Uint8Type, args.UintType,
		args.Float32Type, args.Float64Type:
		return true
	}
	return false
}

func PushArgument(state *lua.State, arg args.Argument) error {
	if arg == nil {
		return nil
	} else if !arg.Valid() {
		return errors.New("cannot push undefined value")
	}

	switch arg.Type() {
	case args.StringType, args.BoolType, args.StringSliceType, args.MapType:
		luar.GoToLua(state, arg.Value())
	case args.NilType:
		state.PushNil()
	case args.CallType:
		luar.GoToLua(state, wrap_call(arg))
	case args.ByteSliceType:
		state.PushBytes(arg.Value().([]byte))
	case args.ArgumentListType, args.ArgumentSliceType:
		return PushArgumentSlice(state, arg, false)
	case args.ArgumentMapType:
		return PushArgumentMap(state, arg.Value().(args.ArgumentMap))
	default:
		if isNumeric(arg) {
			luar.GoToLua(state, arg.Value())
			return nil
		}

		if err := convert_push(arg, state); err != nil {
			if err != ErrCannotConvert {
				return err
			}
		}
		return errors.New("invalid type")
	}

	return nil
}

func PushArgumentSlice(state *lua.State, arg args.Argument, spread bool) error {
	if !arg.Is(args.ArgumentSliceType, args.ArgumentListType) {
		return errors.New("not a slice type")
	}
	top := state.GetTop()
	var val []args.Argument
	if arg.Type() == args.ArgumentSliceType {
		val = arg.Value().([]args.Argument)
	} else {
		val = []args.Argument(arg.Value().(args.ArgumentList))
	}

	if spread {
		for _, a := range val {
			if err := PushArgument(state, a); err != nil {
				state.SetTop(0)
				return err
			}
		}
	} else {

		l := len(val)
		state.CreateTable(l, 0)
		for i := 0; i < l; i++ {
			if err := PushArgument(state, val[i]); err != nil {
				state.SetTop(top)
				return err
			}
			state.RawSeti(-2, i)
		}
	}

	return nil
}

func PushArgumentList(state *lua.State, arg args.ArgumentList) error {
	/*top := state.GetTop()
	l := len(arg)
	state.CreateTable(l, 0)
	for i := 0; i < l; i++ {
		if err := PushArgument(state, arg[i]); err != nil {
			state.SetTop(top)
			return err
		}
		state.RawSeti(-2, i)
	}*/

	return PushArgumentSlice(state, args.Must(arg), true)
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
