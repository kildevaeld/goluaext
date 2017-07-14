//go:generate go-bindata -pkg goluaext -nometadata -o prelude.go prelude.lua

package goluaext

import (
	"bytes"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	yaml "gopkg.in/yaml.v2"

	"github.com/aarzilli/golua/lua"
	uuid "github.com/satori/go.uuid"
	"github.com/stevedonovan/luar"
)

var Globals = []string{
	"pairs",
	"print",
	"rawequal",
	"rawget",
	"rawset",
	"require",
	"select",
	"setfenv",
	"setmetatable",
	"string",
	"table",
	"tonumber",
	"tostring",
	"type",
	"unpack",
	"unsafe_pcall",
	"unsafe_xpcall",
}

type httpOptions struct {
	body []byte
}

func uuidLoader(state *lua.State) int {
	hasher := func(version string) func(state *lua.State) int {

		return func(state *lua.State) int {

			var id uuid.UUID
			switch version {
			case "v1":
				id = uuid.NewV1()
			/*case "v2":
				id = uuid.NewV2()
			case "v3":
				id = uuid.NewV3()*/
			case "v4":
				id = uuid.NewV4()
				/*case "v5":
				id = uuid.NewV5()*/
			}
			//str := state.ToString(1)

			CreateTable(state, luar.Map{
				"bytes": id.Bytes(),
				"valid": id != uuid.Nil,
			}, MetaMap{
				MetaToString: func(state *lua.State) int {
					state.GetField(1, "bytes")
					v := state.ToBytes(-1)
					state.PushString(uuid.FromBytesOrNil(v).String())
					return 1
				},
			})

			return 1
		}
	}

	CreateTable(state, nil, MetaMap{
		MetaIndex: func(state *lua.State) int {
			method := state.ToString(2)
			state.PushGoFunction(hasher(method))
			return 1
		},
	})

	return 1
}

func hashLoader(state *lua.State) int {
	name := state.ToString(1)
	if name != "hash" {
		panic("error")
	}

	hasher := func(algo string) func(state *lua.State) int {

		var h hash.Hash
		switch algo {
		case "sha1":
			h = sha1.New()
		case "sha256":
			h = sha256.New()
		case "sha512":
			h = sha512.New()
		}

		return func(state *lua.State) int {
			if h == nil {
				panic(errors.New("invalid algo"))
			}

			str := state.ToString(1)

			CreateTable(state, luar.Map{
				"bytes": h.Sum([]byte(str)),
			}, MetaMap{
				MetaToString: func(state *lua.State) int {
					state.GetField(1, "bytes")
					v := state.ToBytes(-1)
					state.PushString(fmt.Sprintf("%x", v))
					return 1
				},
			})

			return 1
		}
	}

	CreateTable(state, nil, MetaMap{
		MetaIndex: func(state *lua.State) int {
			method := state.ToString(2)
			state.PushGoFunction(hasher(method))
			return 1
		},
	})

	return 1

}

func getValue(state *lua.State, index int) (interface{}, error) {
	var err error

	var oo interface{}
	t := state.Type(index)
	switch t {
	case lua.LUA_TNIL, lua.LUA_TNONE:
		oo = nil
	case lua.LUA_TTABLE:
		var o luar.Map
		err = luar.LuaToGo(state, index, &o)
		oo = o
	case lua.LUA_TNUMBER:
		oo = state.ToNumber(index)
	case lua.LUA_TSTRING:
		oo = state.ToString(index)
	case lua.LUA_TBOOLEAN:
		oo = state.ToBoolean(index)
	default:
		err = errors.New("invalid type")
	}
	return oo, err
}

func jsonLoader(state *lua.State) int {
	name := state.ToString(1)
	if name != "json" {
		panic("error")
	}

	luar.GoToLua(state, luar.Map{
		"decode": func(str string) (interface{}, error) {
			var out interface{}
			if err := json.Unmarshal([]byte(str), &out); err != nil {
				return nil, err
			}
			return out, nil
		},
		"encode": func(state *lua.State) int {

			var bs []byte
			oo, err := getValue(state, 1)

			indent := ""
			if state.IsString(2) {
				indent = state.ToString(2)
			}

			out := ""
			if err == nil {
				if indent != "" {
					bs, err = json.MarshalIndent(oo, "", indent)
				} else {
					bs, err = json.Marshal(oo)
				}

				if err == nil {
					out = string(bs)
				}
			}
			state.PushString(out)
			if err != nil {
				state.PushString(err.Error())
			} else {
				state.PushNil()
			}

			return 2
		},
	})

	return 1
}

func yamlLoader(state *lua.State) int {

	luar.GoToLua(state, luar.Map{
		"encode": func(state *lua.State) int {
			oo, err := getValue(state, 1)
			var (
				bs  []byte
				out string
			)

			if err == nil {
				bs, err = yaml.Marshal(oo)
				if err == nil {
					out = string(bs)
				}
			}
			state.PushString(out)
			if err != nil {
				state.PushString(err.Error())
			} else {
				state.PushNil()
			}

			return 2
		},
		"decode": func(str string) (interface{}, error) {
			var out interface{}
			if err := yaml.Unmarshal([]byte(str), &out); err != nil {
				return nil, err
			}
			return out, nil
		},
	})

	return 1
}

func httpLoader(state *lua.State) int {
	name := state.ToString(1)
	if name != "http" {
		panic("error")
	}

	request := func(method string) func(state *lua.State) int {
		method = strings.ToUpper(method)

		return func(state *lua.State) int {
			if !state.IsString(1) {
				panic(errors.New("#1 argument must be a string"))
			}

			if state.IsTable(2) {

			}

			var (
				req  *http.Request
				resp *http.Response
				err  error
				bs   []byte
			)

			path := state.ToString(1)
			var options httpOptions
			var body io.Reader

			if options.body != nil {
				body = bytes.NewReader(options.body)
			}

			if req, err = http.NewRequest(method, path, body); err != nil {
				panic(err)
			}

			client := &http.Client{}
			if resp, err = client.Do(req); err != nil {
				//return nil, err.Error()
			}

			if bs, err = ioutil.ReadAll(resp.Body); err != nil {
				//return nil, err.Error()
			}
			fmt.Printf("%s", bs)
			return 0
		}
	}

	state.CreateTable(0, 0)
	state.CreateTable(0, 1)
	state.SetMetaMethod("__index", func(state *lua.State) int {
		method := state.ToString(2)
		state.PushGoFunction(request(method))
		return 1
	})
	state.SetMetaTable(-2)

	return 1
}

func Init() *lua.State {
	state := luar.Init()
	top := state.GetTop()

	state.GetGlobal("package")

	pre := state.GetTop()

	state.CreateTable(0, len(_loaders))
	for k, loader := range _loaders {
		state.PushGoClosure(lua.LuaGoFunction(loader))
		state.SetField(-2, k)
	}

	state.SetField(pre, "preload")
	state.SetTop(top)

	state.MustDoString(string(MustAsset("prelude.lua")))

	return state

}
