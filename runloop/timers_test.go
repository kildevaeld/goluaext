package runloop_test

import (
	"testing"

	"github.com/kildevaeld/goluaext"
	"github.com/kildevaeld/goluaext/runloop"
)

func TestRunloop(t *testing.T) {

	state := goluaext.Init()

	loop, _ := runloop.New(state)

	if err := runloop.RegisterTimers(loop); err != nil {
		t.Fatal(err)
	}

	if err := loop.DoStringAndRun(`
print("Hello from lua")
setTimeout(function()
	print("deferred")
end, 2000)
	`); err != nil {
		t.Fatal(err)
	}

}
