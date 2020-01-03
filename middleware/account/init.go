package account

import (
	"../../session"
)

func InitSession(provider string, addr string, options ...string) (err error) {
	return session.Init(provider, addr, options...)
}
