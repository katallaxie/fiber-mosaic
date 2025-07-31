package fragments

import (
	"sync"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/html"
)

// Resolver ...
type Resolver struct {
	wg sync.WaitGroup
}

// NewResolver ...
func NewResolver() *Resolver {
	return &Resolver{}
}

// ResolverFunc ...
type ResolverFunc func(c *fiber.Ctx, cfg Config) error

// Resolve blocks until all fragments have been called.
//
//nolint:gocyclo
func (r *Resolver) Resolve(c *fiber.Ctx, cfg Config, doc *HTMLFragment) (int, []*html.Node, error) {
	statusCode := fiber.StatusOK
	head := make([]*html.Node, 0)

	ff, err := doc.Fragments()
	if err != nil {
		return statusCode, head, err
	}

	for _, f := range ff {
		r.run(c, cfg, f.Resolve())
	}

	r.wg.Wait()

	for _, f := range ff {
		if f.Primary() && f.statusCode != 0 {
			statusCode = f.statusCode
		}

		fw, err := f.HTMLFragment().Fragments()
		if err != nil {
			return statusCode, head, err
		}

		for _, fwr := range fw {
			if fwr.Ref() == "" {
				continue
			}

			ref, ok := ff[fwr.Ref()]
			if !ok {
				continue
			}

			html, err := ref.HTMLFragment().HTML()
			if err != nil {
				return statusCode, head, err
			}

			fwr.Element().ReplaceWithHtml(html)
		}

		html, err := f.HTMLFragment().HTML()
		if err != nil {
			return statusCode, head, err
		}

		f.Element().ReplaceWithHtml(html)

		head = append(head, f.Links()...)
	}

	return statusCode, head, nil
}

func (r *Resolver) run(c *fiber.Ctx, cfg Config, fn ResolverFunc) {
	r.wg.Add(1)

	go func() {
		defer r.wg.Done()

		err := fn(c, cfg)
		if err != nil {
			return // ignoring errors for now
		}
	}()
}
