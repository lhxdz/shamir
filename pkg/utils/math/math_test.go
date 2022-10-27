package math

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInvMod(t *testing.T) {
	var input, output, mod int64 = 4, 2, 7
	assert.Equal(t, InvMod(input, mod), output)
}
