package values

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPick(t *testing.T) {
	nilled := (*string)(nil)
	empty := ""
	filled := "filled"

	assert.Equal(t, &filled, PickHasValue(nilled, &empty, &filled))
}

func TestStructs(t *testing.T) {
	type A struct {
		B string
	}

	a := A{}
	n := struct{}{}
	assert.True(t, IsZero(a))
	assert.True(t, IsZero(n))
	assert.Equal(t, a, Pick(map[string]A{"a": a}, "a", A{}))
}
