package internal

import (
	"reflect"
	"testing"
)

func AssertEquals(t *testing.T, a, b interface{}, msg string) {
	if a != b {
		t.Fatal(msg)
	}
}

func AssertSuccess(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func AssertFailure(t *testing.T, err error, msg string) {
	if err == nil {
		t.Fatal(msg)
	}
}

func AssertNil(t *testing.T, a interface{}, msg string) {
	if a != nil {
		if reflect.ValueOf(a).Kind() == reflect.Slice {
			if !reflect.ValueOf(a).IsNil() {
				t.Fatal(msg)
			}
		} else {
			t.Fatal(msg)
		}
	}
}

func AssertSlicesEqualU32(t *testing.T, a, b []uint32, msg string) {
	if len(a) != len(b) {
		t.Fatal(msg)
	}

	for i := range a {
		if a[i] != b[i] {
			t.Fatal(msg)
		}
	}
}
