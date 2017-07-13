package goluaext

import (
	"testing"

	args "github.com/kildevaeld/go-args"
	"github.com/stretchr/testify/assert"
)

func TestLua(t *testing.T) {

	state := Init()
	state.OpenLibs()
	if err := state.DoString(`
--util.http.get("http://google.com")
print(util.hash.sha256("Hello, World!"))
print(util.uuid.v3().valid)
function test(msg)
	print("test called", msg)
end

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
