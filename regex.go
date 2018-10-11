package forwardif

import (
	"regexp"

	"github.com/coredns/coredns/request"
)

func NewRegexMatch(regex string) (MatchFunc, error) {
	re, err := regexp.Compile(regex)
	if err != nil {
		return nil, err
	}

	return func(state request.Request) bool {
		return re.MatchString(state.Name())
	}, nil
}
