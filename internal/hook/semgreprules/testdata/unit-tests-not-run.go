//go:build exclude
// +build exclude

package unittestsnotrun

import (
	"testing"

	"github.com/monzo/wearedev/libraries/test"
)

// ruleid:unit-tests-not-run
func TestMain(m *testing.M) {
	dao.Keyspace = gocassa.NewMockKeySpace()
	dao.Init()
}

// ok:unit-tests-not-run
func TestMain(m *testing.M) {
	dao.Keyspace = gocassa.NewMockKeySpace()
	dao.Init()
	test.RunTests(m)
}

// ok:unit-tests-not-run
func TestMain(m *testing.M) {
	dao.Keyspace = gocassa.NewMockKeySpace()
	dao.Init()
	os.Exit(m.Run())
}
