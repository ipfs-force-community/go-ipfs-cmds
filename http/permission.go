package http

import (
	"context"
)

type permKey int

var permCtxKey permKey

func WithPerm(ctx context.Context, perms []string) context.Context {
	return context.WithValue(ctx, permCtxKey, perms)
}

func HasPerm(ctx context.Context, perm string) bool {
	callerPerms, ok := ctx.Value(permCtxKey).([]string)
	if !ok {
		callerPerms = []string{"read"}
	}
	return permContain(perm, callerPerms)
}

type permCode = int

const (
	PermRead  permCode = 0
	PermWrite permCode = 1
	PermSign  permCode = 2
	PermAdmin permCode = 4
)

func convertPermCode(perm string) permCode {
	switch perm {
	case "read":
		return PermRead
	case "write":
		return PermWrite
	case "sign":
		return PermSign
	case "admin":
		return PermAdmin
	default:
		return PermRead
	}
}

func permContain(tagPerm string, callerPerms []string) bool {
	tag := convertPermCode(tagPerm)
	sum := 0
	for _, cp := range callerPerms {
		sum += convertPermCode(cp)
	}
	return tag&sum == tag
}
