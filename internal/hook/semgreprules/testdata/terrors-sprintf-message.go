package tersprintf

import (
	"fmt"

	"github.com/monzo/terrors"
)

func wrap(err error, userID string) error {
	// ruleid:terrors-sprintf-message
	_ = terrors.Augment(err, fmt.Sprintf("process user %s", userID), nil)

	// ruleid:terrors-sprintf-message
	_ = terrors.InternalService("process", fmt.Sprintf("process %s", userID), nil)

	// ok:terrors-sprintf-message
	_ = terrors.Augment(err, "process user", map[string]string{
		"user_id": userID,
	})

	// ok:terrors-sprintf-message
	_ = terrors.BadRequest(fmt.Sprintf("invalid.%s", userID), "process user", nil)

	// ok:terrors-sprintf-message
	_ = terrors.Is(err, fmt.Sprintf("internal_service.%s", userID))

	// ok:terrors-sprintf-message
	return terrors.Augment(err, "process user", nil)
}
