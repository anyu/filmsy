package main

import "testing"

func TestSayHey(t *testing.T) {
	actual := sayHey()
	expected := "hey"

	if actual != expected {
		t.Errorf("got '%s' want '%s'", actual, expected)
	}
}
