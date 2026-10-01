//lint:file-ignore Ignore all the checks
package testdata

import (
	"context"
	"fmt"

	"github.com/monzo/typhon"
)

func example() error {

	// ruleid:request-not-sent
	rawRsp := typhon.NewRequest(
		context.TODO(), "PUT", "/service.foo/bar",
		nil,
	)

	fmt.Println("unrelated")

	if err := rawRsp.Decode(nil); err != nil {
		return nil
	}

	// ok:request-not-sent
	rawFut := typhon.NewRequest(
		context.TODO(), "PUT", "/service.foo/bar",
		nil,
	).Send().Response()

	fmt.Println("unrelated")

	if err := rawFut.Decode(nil); err != nil {
		return nil
	}

	// ok:request-not-sent
	rawRsp = typhon.NewRequest(
		context.TODO(), "PUT", "/service.foo/bar",
		nil,
	)

	fmt.Println("unrelated")

	rawFuture := rawRsp.Send().Response()

	fmt.Println("unrelated")

	if err := rawFuture.Decode(nil); err != nil {
		return nil
	}
	return nil
}
