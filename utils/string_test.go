package utils

import (
	"reflect"
	"testing"
)

func TestJoinStrings(t *testing.T) {
	s := JoinStrings("gzf", "-5", "lf")
	if reflect.DeepEqual(s, "gzf-5lf") == false {
		panic("join strings error")
	}
}
