package forwardif

import (
	"github.com/coredns/coredns/request"
)

type MatchFunc func(request.Request) bool

func FalseMatcher(request.Request) bool {
	return false
}
