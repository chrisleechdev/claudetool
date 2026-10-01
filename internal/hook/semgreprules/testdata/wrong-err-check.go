//lint:file-ignore Ignore all the checks
package testdata

import "fmt"

func handle() {
	err := fmt.Errorf("error")
	var someErr error

	// ruleid:wrong-err-check
	if someErr := something(); err != nil {
		fmt.Println(someErr)
		return
	}

	// ruleid:wrong-err-check
	if someErr = something(); err != nil {
		fmt.Println(someErr)
		return
	}

	// ok:wrong-err-check
	if err := something(); err != nil {
		return
	}

	// ok:wrong-err-check
	if err = something(); err != nil {
		return
	}

	// ok:wrong-err-check
	if _, err := something2(); err != nil {
		return
	}

	// ok:wrong-err-check
	if _, err = something2(); err != nil {
		return
	}

	// ok:wrong-err-check
	if err, _ := something3(); err != nil {
		return
	}

	// ok:wrong-err-check
	if err, _ = something3(); err != nil {
		return
	}

	// ok:wrong-err-check
	if err := nest().Error; err != nil {
		return
	}

	// ok:wrong-err-check
	if err = nest().Error; err != nil {
		return
	}

	// ok:wrong-err-check
	err2 := something()
	if err2 != nil {
		return
	}

	// ok:wrong-err-check
	err = something()
	if err != nil {
		return
	}
}

func something() error {
	return fmt.Errorf("error")
}

func something2() (string, error) {
	return "", nil
}

func something3() (error, bool) {
	return nil, false
}

type Nest struct {
	Error error
}

func nest() Nest {
	return Nest{Error: fmt.Errorf("error")}
}
