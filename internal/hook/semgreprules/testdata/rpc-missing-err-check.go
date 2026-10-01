package main

import (
	"context"

	slackproto "github.com/monzo/wearedev/service.slack/proto"
)

func main() {
	ctx := context.Background()
	err := error(nil)

	// ruleid:rpc-missing-err-check
	slackproto.ReactionsRemoveRequest{}.Send(ctx).DecodeResponse()
	if err != nil {
		return
	}

	req := slackproto.PostRequest{}

	// ruleid:rpc-missing-err-check
	req.Send(ctx).DecodeResponse()
	if err != nil {
		return
	}

	// ok:rpc-missing-err-check
	_, err = slackproto.PostRequest{}.Send(ctx).DecodeResponse()
	if err != nil {
		return
	}

	// Best-effort sends deliberately ignore the error
	// ok:rpc-missing-err-check
	slackproto.PostRequest{}.Send(ctx).DecodeResponse()
}
