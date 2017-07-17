package http

import (
	"testing"

	"github.com/kildevaeld/goluaext"
	"github.com/kildevaeld/goluaext/runloop"
)

func TestHttp(t *testing.T) {

	state := goluaext.Init()

	loop := runloop.New(state)

	if err := RegisterHttp(loop); err != nil {
		t.Fatal(err)
	}

	err := loop.DoStringAndRun(`
local http = require 'http'
print("loaded http")
http.get("https://google.com"):next(function(resp)
	print(resp)
end, function(error)
	print("error",error)
end)
print("made request")
	`)

	if err != nil {
		t.Fatal(err)
	}

}
