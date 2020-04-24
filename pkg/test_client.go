package main

import (
	"context"
	"testing"
)

type clientMock struct{}

var queryMock func(ctx context.Context, q interface{}, variables map[string]interface{}) error

func (c clientMock) Query(ctx context.Context, q interface{}, variables map[string]interface{}) error {
	return queryMock(ctx, &q, variables)
}

func TestQuery(t *testing.T) {
}
