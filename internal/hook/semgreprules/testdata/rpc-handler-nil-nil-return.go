//lint:file-ignore Ignore all the checks
package testdata

import (
	"context"
	"errors"

	"github.com/monzo/terrors"
)

type SomeResponse struct {
	Data string
}

type AnotherResponse struct {
	Value int
}

type Handler struct{}

type SomeRequest struct {
	ID string
}

func someCondition() bool {
	return true
}

func (h *Handler) BadMethod(ctx context.Context, req *SomeRequest) (*SomeResponse, error) {
	if someCondition() {
		// ruleid: rpc-handler-nil-nil-return
		return nil, nil
	}
	return &SomeResponse{Data: "test"}, nil
}

func (h *Handler) AnotherBadMethod(ctx context.Context, req *SomeRequest) (*AnotherResponse, error) {
	// ruleid: rpc-handler-nil-nil-return
	return nil, nil
}

func (h *Handler) GoodMethodWithNonNilError(ctx context.Context, req *SomeRequest) (*SomeResponse, error) {
	// ok: rpc-handler-nil-nil-return
	var err error
	if someCondition() {
		err = errors.New("some error")
	}
	return nil, err
}

func (h *Handler) GoodMethodWithError(ctx context.Context, req *SomeRequest) (*SomeResponse, error) {
	if someCondition() {
		// ok: rpc-handler-nil-nil-return
		return nil, terrors.BadRequest("invalid", "invalid request", nil)
	}
	return &SomeResponse{Data: "test"}, nil
}

func (h *Handler) GoodMethodWithStdError(ctx context.Context, req *SomeRequest) (*SomeResponse, error) {
	if someCondition() {
		// ok: rpc-handler-nil-nil-return
		return nil, errors.New("some error")
	}
	return &SomeResponse{Data: "test"}, nil
}

func (h *Handler) GoodMethodWithResponse(ctx context.Context, req *SomeRequest) (*SomeResponse, error) {
	// ok: rpc-handler-nil-nil-return
	return &SomeResponse{Data: "test"}, nil
}

// Not an RPC handler: doesn't return *Response
func (h *Handler) NotAnRPCHandler(ctx context.Context, req *SomeRequest) (string, error) {
	// ok: rpc-handler-nil-nil-return
	return "", nil
}

// Not a handler method
func regularFunction() (*SomeResponse, error) {
	// ok: rpc-handler-nil-nil-return
	return nil, nil
}
