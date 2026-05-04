package validation

import (
	"fmt"
	"reflect"
	"testing"
)

type TestStruct struct {
	foo int64
	bar string
}

func TestIsZero(t *testing.T) {
	tests := []TestStruct{
		{0, ""},
		{4323423423423423423, "a"},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt), func(t *testing.T) {
			if ok := IsStructFieldEmpty(reflect.ValueOf(tt.foo)); !ok {
				t.Errorf("IsStructFieldEmpty() = %v, want %v", ok, tt.foo)
			}
		})
	}
}
