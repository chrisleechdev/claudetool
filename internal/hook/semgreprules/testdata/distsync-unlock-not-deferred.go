package main

import (
	"context"

	"github.com/monzo/wearedev/libraries/distsync"
)

func main() {
	lock1, err := distsync.Lock(context.Background(), "some-lock")
	if err != nil {
		return
	}
	// ok:distsync-unlock-not-deferred
	defer lock1.Unlock()

	lock2, err := distsync.Lock(context.Background(), "some-lock")
	if err != nil {
		return
	}
	// ruleid:distsync-unlock-not-deferred
	lock2.Unlock()
}
