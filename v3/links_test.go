package fragments_test

import (
	"testing"

	fragments "github.com/katallaxie/fiber-mosaic/v3"

	"github.com/stretchr/testify/assert"
)

func TestHeader_Links(t *testing.T) {
	tests := []struct {
		desc string
		in   fragments.Header
		want []fragments.Link
	}{
		{
			desc: "should return a list of links",
			in:   fragments.Header("<https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css>; rel=\"stylesheet\", <https://unpkg.com/react-dom@17/umd/react-dom.development.js>; rel=\"script\""),
			want: []fragments.Link{
				{
					URL:    "https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css",
					Rel:    "stylesheet",
					Params: map[string]string{},
				},
				{
					URL:    "https://unpkg.com/react-dom@17/umd/react-dom.development.js",
					Rel:    "script",
					Params: map[string]string{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := test.in.Links()
			assert.Equal(t, test.want, got)
		})
	}
}

func TestHeader_FilterByStylesheet(t *testing.T) {
	tests := []struct {
		desc string
		in   fragments.Header
		want []fragments.Link
	}{
		{
			desc: "should return a list of links",
			in:   fragments.Header("<https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css>; rel=\"stylesheet\", <https://unpkg.com/react-dom@17/umd/react-dom.development.js>; rel=\"script\""),
			want: []fragments.Link{
				{
					URL:    "https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css",
					Rel:    "stylesheet",
					Params: map[string]string{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := fragments.FilterByStylesheet(test.in.Links()...)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestHeader_FilterByScript(t *testing.T) {
	tests := []struct {
		desc string
		in   fragments.Header
		want []fragments.Link
	}{
		{
			desc: "should return a list of links",
			in:   fragments.Header("<https://unpkg.com/tailwindcss@^2/dist/tailwind.min.css>; rel=\"stylesheet\", <https://unpkg.com/react-dom@17/umd/react-dom.development.js>; rel=\"script\""),
			want: []fragments.Link{
				{
					URL:    "https://unpkg.com/react-dom@17/umd/react-dom.development.js",
					Rel:    "script",
					Params: map[string]string{},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			got := fragments.FilterByScript(test.in.Links()...)
			assert.Equal(t, test.want, got)
		})
	}
}
