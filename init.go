package goluaext

import args "github.com/kildevaeld/go-args"

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
