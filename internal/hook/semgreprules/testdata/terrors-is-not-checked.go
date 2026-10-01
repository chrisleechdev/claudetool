package terrorsisnotchecked

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/monzo/terrors"

	"github.com/monzo/wearedev/libraries/test"
)

func Test(t *testing.T) {
	err := terrors.InternalService("oops", "Oops!", map[string]string{
		"foo": "bar",
	})

	// ruleid:terrors-is-not-checked
	terrors.Is(err, terrors.ErrInternalService, "oops")

	// ok:terrors-is-not-checked
	assert.True(t, terrors.Is(err, terrors.ErrInternalService, "oops"))

	// ok:terrors-is-not-checked
	assert.False(t, terrors.Is(err, terrors.ErrInternalService, "oops"), "Expected an oops")

	// ok:terrors-is-not-checked
	require.True(t, terrors.Is(err, terrors.ErrInternalService, "oops"))

	// ok:terrors-is-not-checked
	require.False(t, terrors.Is(err, terrors.ErrInternalService, "oops"), "Required an oops")

	// ok:terrors-is-not-checked
	ok := terrors.Is(err, terrors.ErrInternalService, "oops")
	assert.True(t, ok)

	// ok:terrors-is-not-checked
	ok = terrors.Is(err, terrors.ErrInternalService, "oops")
	assert.True(t, ok)

	// ok:terrors-is-not-checked
	ok2 := terrors.Is(err, terrors.ErrInternalService, "oops")
	require.True(t, ok2)

	// ok:terrors-is-not-checked
	ok2 = terrors.Is(err, terrors.ErrInternalService, "oops")
	require.True(t, ok2)

	// ok:terrors-is-not-checked
	if !terrors.Is(err, terrors.ErrInternalService, "oops") {
		t.Errorf("Expected an oops")
	}

	// ok:terrors-is-not-checked
	switch {
	case false,
		terrors.Is(err, terrors.ErrInternalService, "oops"):
		t.Errorf("Expected an oops")
	}

	// ok:terrors-is-not-checked
	test.AssertTerrorMatches(t, err, terrors.ErrInternalService, "oops")
}
