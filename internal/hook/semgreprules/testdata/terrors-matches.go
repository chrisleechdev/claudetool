package termatches

import "github.com/monzo/terrors"

func check(err error) bool {
	// ruleid:terrors-matches
	if terrors.Matches(err, "bad_request.foo") {
		return true
	}

	// ruleid:terrors-matches
	if terrors.PrefixMatches(err, "bad_request") {
		return true
	}

	// ok:terrors-matches
	return terrors.Is(err, terrors.ErrBadRequest, "foo")
}
