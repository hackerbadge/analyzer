package main

import "testing"

func TestAppendUnique(t *testing.T) {
	var a = []string{"foo", "bar"}
	a = AppendUnique(a, "hehe")
	a = AppendUnique(a, "foo")
	a = AppendUnique(a, "hehe")
	a = AppendUnique(a, "bar")
	a = AppendUnique(a, "blah")
}
