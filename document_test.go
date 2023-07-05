package fragments_test

import (
	"testing"

	fragments "github.com/katallaxie/fiber-mosaic"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/html"
)

func TestNewDocument(t *testing.T) {
	var tests = []struct {
		desc string
		root *html.Node
	}{
		{
			desc: "should create a new document",
			root: &html.Node{
				Type:     html.DocumentNode,
				DataAtom: 0,
				Data:     "",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			doc, err := fragments.NewDocument(test.root)
			assert.NoError(t, err)
			assert.NotNil(t, doc)
		})
	}
}
