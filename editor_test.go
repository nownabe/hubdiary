package main

import "testing"

func Test_editor(t *testing.T) {
	e := editor{"test/editor.sh"}
	actual, err := e.Edit("")

	if err != nil {
		t.Fatalf("error is not expected: %v", err)
	}

	expect := "this line was added by editor.sh"

	if actual != expect {
		t.Errorf("Expected '%s', but got '%s'", expect, actual)
	}
}
