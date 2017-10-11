package goluaext

import (
	"errors"
	"fmt"

	"github.com/aarzilli/golua/lua"
	args "github.com/kildevaeld/go-args"
)

var _converters map[args.Type]Converter
var _loaders map[string]LuaCallback

func init() {
	_converters = make(map[args.Type]Converter)
	_loaders = make(map[string]LuaCallback)

	_loaders["json"] = jsonLoader
	_loaders["http"] = httpLoader
	_loaders["hash"] = hashLoader
	_loaders["uuid"] = uuidLoader
	_loaders["yaml"] = yamlLoader
}

func RegisterModule(name string, fn LuaCallback, overwrite bool) error {
	if _, ok := _loaders[name]; ok && !overwrite {
		return errors.New("a module with that name already exists")
	}

	_loaders[name] = fn

	return nil
}

func RegisterModuleOnVM(state *lua.State, name string, fn LuaCallback, overwrite bool) error {
	if _, ok := _loaders[name]; ok && !overwrite {
		return errors.New("a module with that name already exists")
	}

	top := state.GetTop()

	state.GetGlobal("package")

	pre := state.GetTop()

	state.GetField(-1, "preload")
	if !state.IsTable(-1) {
		state.CreateTable(0, 1)
	}

	state.PushGoClosure(lua.LuaGoFunction(fn))
	state.SetField(-2, name)

	state.SetField(pre, "preload")
	state.SetTop(top)

	return nil
}

func RegisterLuaModuleOnVM(state *lua.State, name string, luaString string, overwrite bool) error {
	if _, ok := _loaders[name]; ok && !overwrite {
		return errors.New("a module with that name already exists")
	}

	top := state.GetTop()

	state.GetGlobal("package")

	pre := state.GetTop()

	state.GetField(-1, "preload")
	if !state.IsTable(-1) {
		state.CreateTable(0, 1)
	}

	state.PushGoClosure(func(state *lua.State) int {

		if err := state.DoString(luaString); err != nil {
			panic(err)
		}

		if state.IsTable(-1) {
			//fmt.Printf("it's table")
		}

		return 1
	})
	state.SetField(-2, name)

	state.SetField(pre, "preload")
	state.SetTop(top)

	return nil
}

func RegisterLuaModule(state *lua.State, name string, luaString string, overwrite bool) error {
	if _, ok := _loaders[name]; ok && !overwrite {
		return errors.New("a module with that name already exists")
	}

	_loaders[name] = func(state *lua.State) int {

		if err := state.DoString(luaString); err != nil {
			panic(err)
		}

		if state.IsTable(-1) {
			//fmt.Printf("it's table")
		}

		return 1
	}

	return nil
}
