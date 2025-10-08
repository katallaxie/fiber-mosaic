package fragments_test

import (
	"testing"

	fragments "github.com/katallaxie/fiber-mosaic/v3"

	"github.com/stretchr/testify/assert"
)

func TestNewClient(t *testing.T) {
	client := fragments.NewClient()

	assert.NotNil(t, client)
}
