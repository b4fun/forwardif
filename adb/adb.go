package adb

import (
	"io"
	"os"
	"strings"

	"github.com/b4fun/adblockdomain"
	"github.com/coredns/coredns/request"
)

type AdBlock struct {
	domains []string
}

func NewAdBlock(source io.Reader, useException bool) (*AdBlock, error) {
	f := adblockdomain.ParseFromReader
	if useException {
		f = adblockdomain.ParseExceptionFromReader
	}
	domains, err := f(source)
	if err != nil {
		return nil, err
	}

	return &AdBlock{domains: domains}, nil
}

func NewAdBlockFromFilePath(path string, useException bool) (*AdBlock, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewAdBlock(file, useException)
}

func (a AdBlock) Match(state request.Request) bool {
	name := strings.TrimSuffix(state.Name(), ".")

	// TODO: speed
	for _, domain := range a.domains {
		if name == domain {
			return true
		}
		if strings.HasSuffix(name, domain) {
			return true
		}
	}

	return false
}
