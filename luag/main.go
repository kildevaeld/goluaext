package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/kildevaeld/goluaext"
)

func main() {

	a := os.Args[1:]

	if len(a) == 0 {
		fmt.Printf("usage: luag <path>")
		os.Exit(1)
	}
	path := a[0]
	if !filepath.IsAbs(path) {
		var err error
		if path, err = filepath.Abs(path); err != nil {
			fmt.Fprintf(os.Stderr, "not a path\n")
			os.Exit(1)
		}
	}
	oldLuaPath := os.Getenv("LUA_PATH")
	luaPath := oldLuaPath
	if luaPath != "" {
		luaPath += ";"
	}

	luaPath += filepath.Dir(path) + "/?.lua"

	os.Setenv("LUA_PATH", luaPath)
	lua := goluaext.Init()
	os.Setenv("LUA_PATH", oldLuaPath)

	if err := lua.DoFile(path); err != nil {
		fmt.Fprintf(os.Stderr, "error %s\n", err)
		os.Exit(1)
	}

}
