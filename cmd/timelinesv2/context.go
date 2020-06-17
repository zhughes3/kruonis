package main

import "context"

type accessTokenClaimsContextKey struct{}
type refreshTokenClaimsContextKey struct{}
type userIDContextKey struct{}

// AccessTokenClaimsFromContext returns the access token claims stored in a context, or nil if there isn't one.
func AccessTokenClaimsFromContext(ctx context.Context) accessTokenClaims {
	claims, _ := ctx.Value(accessTokenClaimsContextKey{}).(accessTokenClaims)
	return claims
}

// NewContextWithAccessTokenClaims returns a new context with access token claims.
func NewContextWithAccessTokenClaims(parent context.Context, c accessTokenClaims) context.Context {
	return context.WithValue(parent, accessTokenClaimsContextKey{}, c)
}

// RefreshTokenClaimsFromContext returns the refresh token claims stored in a context, or nil if there isn't one.
func RefreshTokenClaimsFromContext(ctx context.Context) refreshTokenClaims {
	claims, _ := ctx.Value(refreshTokenClaimsContextKey{}).(refreshTokenClaims)
	return claims
}

// NewContextWithRefreshTokenClaims returns a new context with refresh token claims.
func NewContextWithRefreshTokenClaims(parent context.Context, c refreshTokenClaims) context.Context {
	return context.WithValue(parent, refreshTokenClaimsContextKey{}, c)
}
