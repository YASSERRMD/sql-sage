package models

import (
	"testing"
)

func TestObjectTypeValid(t *testing.T) {
	cases := []struct {
		in   ObjectType
		want bool
	}{
		{ObjectProcedure, true},
		{ObjectFunction, true},
		{ObjectPackage, true},
		{ObjectTrigger, true},
		{ObjectView, true},
		{ObjectScript, true},
		{ObjectUnknown, true},
		{ObjectType("bogus"), false},
	}
	for _, c := range cases {
		if got := c.in.Valid(); got != c.want {
			t.Errorf("%s: want %v got %v", c.in, c.want, got)
		}
	}
}
