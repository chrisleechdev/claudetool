package terrortestmatching

import (
	"testing"

	"github.com/monzo/terrors"
	"github.com/monzo/wearedev/libraries/test"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	err := terrors.InternalService("ohno", "Oh no!", nil)

	// ruleid:terror-test-matching
	assert.True(t, terrors.Is(err, "internal_service.ohno"))

	// ruleid:terror-test-matching
	assert.True(t, terrors.Is(err, terrors.ErrInternalService, "ohno"))

	// ruleid:terror-test-matching
	assert.True(t, terrors.Is(err, terrors.ErrInternalService, "ohno"), "Expected an ohno")

	// ok:terror-test-matching
	test.AssertTerrorMatches(t, err, "internal_service.ohno")
}
