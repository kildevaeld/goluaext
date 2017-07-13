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

	"github.com/aarzilli/golua/lua"
	"github.com/stevedonovan/luar"
)

type httpOptions struct {
	body []byte
}

func createHttp(state *lua.State) {
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
	top := state.GetTop()

	state.GetGlobal("util")
	state.PushValue(top)
	state.SetField(-2, "http")
}

func createHash(state *lua.State) {

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

			state.PushBytes(h.Sum([]byte(str)))

			return 1
		}
	}

	state.CreateTable(0, 0)
	state.CreateTable(0, 1)
	state.SetMetaMethod("__index", func(state *lua.State) int {
		method := state.ToString(2)
		state.PushGoFunction(hasher(method))
		return 1
	})
	state.SetMetaTable(-2)
	top := state.GetTop()

	state.GetGlobal("util")
	state.PushValue(top)
	state.SetField(-2, "hash")
}

func Init() *lua.State {
	state := luar.Init()

	luar.Register(state, "util", luar.Map{
		"json": luar.Map{
			"decode": func(str string) (luar.Map, error) {
				var out luar.Map
				if err := json.Unmarshal([]byte(str), &out); err != nil {
					return nil, err
				}
				return out, nil
			},
			"encode": func(state *lua.State) int {

				var err error
				var bs []byte
				var oo interface{}
				t := state.Type(1)
				switch t {
				case lua.LUA_TNIL, lua.LUA_TNONE:
					oo = nil
				case lua.LUA_TTABLE:
					var o luar.Map
					err = luar.LuaToGo(state, 1, &o)
					oo = o
				case lua.LUA_TNUMBER:
					oo = state.ToNumber(1)
				case lua.LUA_TSTRING:
					oo = state.ToString(1)
				case lua.LUA_TBOOLEAN:
					oo = state.ToBoolean(1)
				default:
					err = errors.New("invalid type")
				}

				out := ""
				if err == nil {
					bs, err = json.Marshal(oo)
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
		},
	})

	createHttp(state)
	createHash(state)

	state.MustDoString(string(MustAsset("prelude.lua")))

	return state

}
