package main

import "context"

type userIDContextKey struct{}

// UserIDFromContext returns the ID stored in a context, or nil if there isn't one.
func UserIDFromContext(ctx context.Context) uint64 {
	id, _ := ctx.Value(userIDContextKey{}).(uint64)
	return id
}

// NewContextWithUserID returns a new context with the given Span attached.
func NewContextWithUserID(parent context.Context, id uint64) context.Context {
	return context.WithValue(parent, userIDContextKey{}, id)
}
