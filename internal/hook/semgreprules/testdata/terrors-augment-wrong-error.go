//lint:file-ignore Ignore all the checks
package testdata

import (
	"errors"

	"github.com/monzo/terrors"
)

func tup() (string, error) {
	if []bool{true}[0] {
		return "", errors.New("err")
	} else {
		return "foo", nil
	}
}

func terrorsAugmentFine() error {
	_, e := tup()

	if e != nil {
		// ok:terrors-augment-wrong-error
		return terrors.Augment(e, "some message", nil)
	}

	return nil
}

func terrorsAugmentTupFine() (string, error) {
	v, e := tup()

	// ok:terrors-augment-wrong-error
	if e != nil {
		return "", terrors.Augment(e, "some message", nil)
	}

	return v, nil
}

func nestedIf() (string, error) {
	var b *bool

	if b != nil {
		v, err := tup()
		if err != nil {
			// ok:terrors-augment-wrong-error
			return "", terrors.Augment(err, "error getting thing", nil)
		}
		return v, nil
	}

	return "", nil
}

func terrorsAugment() error {
	var err error

	_, e := tup()

	if e != nil {
		// ruleid:terrors-augment-wrong-error
		return terrors.Augment(err, "some message", nil)
	}

	return nil
}

func terrorsAugmentTup() (string, error) {
	var err error

	v, e := tup()

	if e != nil {
		// ruleid:terrors-augment-wrong-error
		return "", terrors.Augment(err, "some message", nil)
	}

	return v, nil
}

func shadowedError() error {
	_, err := tup()
	if err != nil {
		if _, err := tup(); err != nil {
			// ok:terrors-augment-wrong-error
			return terrors.Augment(err, "error doing inner thing", nil)
		}

		// ok:terrors-augment-wrong-error
		return terrors.Augment(err, "error doing thing", nil)
	}

	return nil
}

func errorSwitch() error {
	var b *bool
	v, err := tup()
	if b != nil {
		switch {
		case err != nil:
			// ok:terrors-augment-wrong-error
			return terrors.Augment(err, "error doing thing", nil)
		case v != "":
			return nil
		default:
			return nil
		}
	}

	return nil
}

func errorCreatedWithinBranch() error {
	var b *bool
	if b != nil {
		err := errors.New("err")
		// ok:terrors-augment-wrong-error
		return terrors.Augment(err, "error doing thing", nil)
	}

	return nil
}
