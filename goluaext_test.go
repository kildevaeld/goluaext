package goluaext

import (
	"testing"

	args "github.com/kildevaeld/go-args"
	"github.com/stretchr/testify/assert"
)

func TestLua(t *testing.T) {

	state := Init()
	//state.OpenLibs()
	if err := state.DoString(`
--util.http.get("http://google.com")
local hash = require 'hash'
local uuid = require 'uuid'
print(hash.sha256("Hello, World!"))
print(uuid.v4())
function test(msg)
	print("test called", msg)
end

local json = require 'json'
local yaml = require 'yaml'

local js = json.encode({
	test = "hello json"
})

print(js)

local l = json.decode(js)

print(l.test)

print(yaml.encode({
	test = "hello yaml"
}))

local http = require 'http'

js = json.encode("rapper")

print(json.decode(js))

http.get("https://google.com")


	`); err != nil {
		t.Fatal(err)
	}

	state.GetGlobal("test")
	assert.True(t, state.IsFunction(-1))
	a, e := LuaToArgument(state, -1, true)
	if e != nil {
		t.Fatal(e)
	}
	defer a.Free()

	a.(args.Call).Call(args.ArgumentList{args.Must("World")})

}
