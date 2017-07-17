package http

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aarzilli/golua/lua"
	"github.com/kildevaeld/goluaext"
	"github.com/kildevaeld/goluaext/runloop"
	"github.com/stevedonovan/luar"
)

type httpOptions struct {
	body    []byte
	url     string
	method  string
	headers http.Header
}

type HttpRequest struct {
	Body    []byte
	Headers map[string]string
}

func doRequest(options httpOptions) (*http.Response, error) {

	var (
		req  *http.Request
		err  error
		body io.Reader
	)

	if options.body != nil {
		body = bytes.NewReader(options.body)
	}

	if req, err = http.NewRequest(options.method, options.url, body); err != nil {
		return nil, err
	}

	if options.headers != nil {
		req.Header = options.headers
	}

	client := &http.Client{}
	return client.Do(req)
}

func InitPromise(state *lua.State) error {
	if err := goluaext.Require(state, "promise"); err != nil {
		return err
	}

	state.GetField(-1, "new")
	if !state.IsFunction(-1) {
		return errors.New("invalid promise module")
	}

	if err := state.Call(0, 1); err != nil {
		return err
	}

	if !state.IsTable(-1) {
		return errors.New("invalid promise")
	}
	return nil
}

func requestFactory(method string, loop *runloop.Loop) lua.LuaGoFunction {
	method = strings.ToUpper(method)

	return func(state *lua.State) int {
		if !state.IsString(1) {
			panic(errors.New("#1 argument must be a string"))
		}

		o := httpOptions{
			url: state.ToString(1),
		}

		if state.IsTable(2) {
			var request HttpRequest
			if err := luar.LuaToGo(state, 2, &request); err != nil {
				panic(err)
			}
			o.body = request.Body

			for k, v := range request.Headers {
				o.headers.Set(k, v)
			}
		}

		if err := InitPromise(state); err != nil {
			panic(err)
		}

		promise := luar.NewLuaObject(state, -1)

		task := &http_task{
			p: promise,
		}

		loop.Add(task)

		go func() {
			defer loop.Ready(task)
			resp, err := doRequest(o)
			//fmt.Printf("rapper %s", err)
			if err != nil {

				task.err = err
				return
			}

			defer resp.Body.Close()

			bs, err := ioutil.ReadAll(resp.Body)
			if err != nil {

				task.err = err

			} else {

				task.content = luar.Map{
					"status": resp.StatusCode,
					"body":   bs,
				}
			}

		}()

		return 1
	}
}

func RegisterHttp(loop *runloop.Loop) error {

	return goluaext.RegisterModuleOnVM(loop.VM(), "http", func(state *lua.State) int {

		goluaext.CreateTable(state, nil, goluaext.MetaMap{
			goluaext.MetaIndex: func(state *lua.State) int {
				name := state.ToString(2)
				state.PushGoFunction(requestFactory(name, loop))
				return 1
			},
		})
		return 1
	}, true)

}
