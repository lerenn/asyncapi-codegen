package asyncapiv3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFollowReturnsSelfWithoutReference ensures Follow returns the structure
// itself when it does not reference another one (issue #146).
func TestFollowReturnsSelfWithoutReference(t *testing.T) {
	srv := &Server{}
	assert.Same(t, srv, srv.Follow())

	tag := &Tag{}
	assert.Same(t, tag, tag.Follow())

	ex := &MessageExample{}
	assert.Same(t, ex, ex.Follow())

	ora := &OperationReplyAddress{}
	assert.Same(t, ora, ora.Follow())

	sec := &SecurityScheme{}
	assert.Same(t, sec, sec.Follow())

	sv := &ServerVariable{}
	assert.Same(t, sv, sv.Follow())

	doc := &ExternalDocumentation{}
	assert.Same(t, doc, doc.Follow())
}

// TestFollowReturnsReference ensures Follow returns the referenced structure
// when one is set (issue #146).
func TestFollowReturnsReference(t *testing.T) {
	target := &Server{}
	ref := &Server{ReferenceTo: target}
	assert.Same(t, target, ref.Follow())

	tagTarget := &Tag{}
	tagRef := &Tag{ReferenceTo: tagTarget}
	assert.Same(t, tagTarget, tagRef.Follow())

	secTarget := &SecurityScheme{}
	secRef := &SecurityScheme{ReferenceTo: secTarget}
	assert.Same(t, secTarget, secRef.Follow())
}
