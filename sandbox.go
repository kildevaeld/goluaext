package goluaext

import (
	"errors"

	"github.com/aarzilli/golua/lua"
)

func CreateSandbox(state *lua.State, meta MetaMap, globals []string) error {
	top := state.GetTop()
	if !state.IsFunction(top) {
		return errors.New("not a function")
	}

	state.CreateTable(0, len(globals))
	for _, g := range globals {
		state.GetGlobal(g)
		state.SetField(-2, g)
	}

	state.CreateTable(0, len(meta))
	for k, m := range meta {
		state.SetMetaMethod(k, m)
	}
	state.SetMetaTable(-2)

	state.SetfEnv(top)

	return nil
}
