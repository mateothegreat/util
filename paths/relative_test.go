package paths

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindProjectRoot(t *testing.T) {
	root, err := FindRootByMarker("go.mod")
	if err != nil {
		t.Fatal(err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, fmt.Sprintf("%s/%s", root, path.Base(wd)), wd)
}
