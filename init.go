package goluaext

import (
	"errors"

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
