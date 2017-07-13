package goluaext

import "testing"

func TestLua(t *testing.T) {

	state := Init()
	state.OpenLibs()
	if err := state.DoString(`
--util.http.get("http://google.com")
print(util.hash.sha256("Hello, World!"))
	`); err != nil {
		t.Fatal(err)
	}

}
