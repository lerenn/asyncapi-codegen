package asyncapiv2

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestParameterFollow ensures Parameter.Follow returns the referenced parameter
// when set and the parameter itself otherwise (issue #146).
func TestParameterFollow(t *testing.T) {
	p := &Parameter{}
	assert.Same(t, p, p.Follow())

	target := &Parameter{}
	ref := &Parameter{ReferenceTo: target}
	assert.Same(t, target, ref.Follow())
}
