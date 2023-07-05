package fragments_test

import (
	"testing"

	fragments "github.com/katallaxie/fiber-mosaic"
	"github.com/stretchr/testify/assert"
)

func TestNewResolver(t *testing.T) {
	var tests = []struct {
		desc string
	}{}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			r := fragments.NewResolver()
			assert.NotNil(t, r)
		})
	}
}
