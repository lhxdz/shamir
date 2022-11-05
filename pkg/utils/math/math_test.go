package math

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInvMod(t *testing.T) {
	var input, output, mod int64 = 4, 2, 7
	assert.Equal(t, InvMod(input, mod), output)
}
