package nouseafternilcheck

import "github.com/monzo/terrors"

func somethingThatReturnsError() error                    { return nil }
func somethingThatReturnsMultipleValues() (string, error) { return "", nil }
func someCondition() bool                                 { return true }

func useAfterNil() error {
	err := somethingThatReturnsError()
	// ruleid:no-use-after-nil-check
	if err != nil {
		return err
	}
	return err
}

// Only err-named variables are checked
func allowsOtherVariables() error {
	something := somethingThatReturnsError()
	if something != nil {
		return something
	}
	// ok:no-use-after-nil-check
	return something
}

func augmentAfterNil() error {
	err := somethingThatReturnsError()
	// ruleid:no-use-after-nil-check
	if err != nil {
		return err
	}
	if someCondition() {
		return terrors.Augment(err, "some context", nil)
	}
	return nil
}

func redefinedAfterNilCheck() error {
	err := somethingThatReturnsError()
	if err != nil {
		return err
	}
	// ok:no-use-after-nil-check
	err = somethingThatReturnsError()
	// ok:no-use-after-nil-check
	return err
}

func useAfterNilMultipleValues() (string, error) {
	someString, err := somethingThatReturnsMultipleValues()
	// ruleid:no-use-after-nil-check
	if err != nil {
		return "", err
	}
	return someString, err
}

func redefineAfterNilMultipleValues() (string, error) {
	someString, err := somethingThatReturnsMultipleValues()
	if err != nil {
		return "", err
	}
	// ok:no-use-after-nil-check
	anotherString, err := somethingThatReturnsMultipleValues()
	// ok:no-use-after-nil-check
	return someString + anotherString, err
}
