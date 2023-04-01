package multihasher_test

import (
	"strings"
	"testing"

	"github.com/overlordtm/pmss/pkg/multihasher"
)

func TestMultiHash(t *testing.T) {
	h, err := multihasher.Hash(strings.NewReader(`Lorem ipsum dolor sit amet`))

	if err != nil {
		t.Error(err)
	}

	if h.MD5 != "fea80f2db003d4ebc4536023814aa885" {
		t.Error("MD5 hash mismatch")
	}

	if h.SHA1 != "38f00f8738e241daea6f37f6f55ae8414d7b0219" {
		t.Error("SHA1 hash mismatch")
	}

	if h.SHA256 != "16aba5393ad72c0041f5600ad3c2c52ec437a2f0c7fc08fadfc3c0fe9641d7a3" {
		t.Error("SHA256 hash mismatch")
	}
}
