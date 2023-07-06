package fragments_test

import (
	"testing"

	fragments "github.com/katallaxie/fiber-mosaic"

	"github.com/stretchr/testify/assert"
)

func TestNewClei(t *testing.T) {
	client := fragments.NewClient()

	assert.NotNil(t, client)
}
