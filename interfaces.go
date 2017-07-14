package goluaext

import "github.com/aarzilli/golua/lua"

var (
	MetaIndex    = "__index"
	MetaNewIndex = "__newindex"
	MetaToString = "__tostring"
	MetaCall     = "__call"
	MetaLen      = "__len"
	MetaGC       = "__gc"
	MetaEqual    = "__eq"
	MetaLT       = "__lt"
	MetaLE       = "__le"
	MetaAdd      = "__add"
	MetaSub      = "__sub"
	MetaMul      = "__mul"
	MetaDiv      = "__div"
	MetaMod      = "__mod"
	MetaPow      = "__pow"
	MetaConcat   = "__concat"
	//MetaMeta     = "__meta"
)

type MetaMap map[string]func(*lua.State) int

type LuaCallback func(state *lua.State) int
