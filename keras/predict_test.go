package main

import (
	"testing"

	"github.com/matryer/is"
)

func TestBodyToVector(t *testing.T) {
	is := is.New(t)
	body := "Here is some text. I Love Go"

	wordIndex := map[string]int32{
		"here": 1,
		"go":   2,
		"text": 3,
		"love": 4,
		"I":    5,
	}

	x := bodyToVector(body, wordIndex, 20)

	// x = [ 0 0 0 0 0 0 0 0 0 0, 1, 0, 0, 3, 5, 4, 2]
	// check tokens
	is.Equal(len(x), 20)
	t.Log(x)
	is.Equal(x[17], wordIndex["i"])    // I
	is.Equal(x[18], wordIndex["love"]) // Love
	is.Equal(x[19], wordIndex["go"])   // Go
}
