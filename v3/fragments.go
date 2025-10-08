// 🚀 Fiber is an Express inspired web framework written in Go with 💖
// 📌 API Documentation: https://fiber.wiki
// 📝 Github Repository: https://github.com/gofiber/fiber

package fragments

import (
	"bytes"
	"html/template"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v3"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/html"
)

// DefaultClient is the default http client used.
var DefaultClient = NewClient()

// Config ...
type Config struct {
	// Filter defines a function to skip the middleware.
	// Optional. Default: nil
	Filter func(fiber.Ctx) bool

	// FilterResponse defines a function to filter the responses
	// from the fragment sources.
	FilterResponse func(*fasthttp.Response)

	// FilterRequest defines a function to filter the request
	// to the fragment sources.
	FilterRequest func(*fasthttp.Request)

	// ErrorHandler defines a function which is executed
	// It may be used to define a custom error.
	// Optional. Default: 401 Invalid or expired key
	ErrorHandler fiber.ErrorHandler

	// FilterHead defines a function to filter the new
	// nodes in the <head> of the document passed by the LINK header entity.
	FilterHead func([]*html.Node) []*html.Node

	// DefaultHost defines the host to use,
	// if no host is set on a fragment.
	// Optional. Default: localhost:3000
	DefaultHost string

	// Client defines the http client to use for
	// fetching the fragments.
	// Optional. Default: DefaultClient
	Client *http.Client
}

// Template is a middleware for templating.
//
//nolint:gocyclo
func Template(config Config, name string, bind interface{}, layouts ...string) fiber.Handler {
	// Set default config
	cfg := configDefault(config)

	return func(c fiber.Ctx) error {
		// Filter request to skip middleware
		if cfg.Filter != nil && cfg.Filter(c) {
			return c.Next()
		}

		method := c.Method()

		// We only handle GET and HEAD requests
		if method != fiber.MethodGet && method != fiber.MethodHead {
			return c.Next()
		}

		var err error
		buf := new(bytes.Buffer)

		// Render raw template using 'name' as filepath if no engine is set
		var tmpl *template.Template
		if _, err = readContent(buf, name); err != nil {
			return cfg.ErrorHandler(c, err)
		}
		// Parse template
		if tmpl, err = template.New("").Parse(buf.String()); err != nil {
			return cfg.ErrorHandler(c, err)
		}
		buf.Reset()
		// Render template
		if err = tmpl.Execute(buf, bind); err != nil {
			return cfg.ErrorHandler(c, err)
		}

		if c.App().Config().Views != nil {
			// Render template based on global layout if exists
			if len(layouts) == 0 && c.App().Config().ViewsLayout != "" {
				layouts = []string{
					c.App().Config().ViewsLayout,
				}
			}
			// Render template from Views
			if err := c.App().Config().Views.Render(buf, name, bind, layouts...); err != nil {
				return cfg.ErrorHandler(c, err)
			}
		}

		r := bytes.NewReader(buf.Bytes())

		root, e := html.Parse(r)
		if e != nil {
			return cfg.ErrorHandler(c, err)
		}

		doc, err := NewDocument(root)
		if err != nil {
			return cfg.ErrorHandler(c, err)
		}

		return Do(c, cfg, doc)
	}
}

// Do represents the core functionality of the middleware.
// It resolves the fragments from a parsed template.
func Do(c fiber.Ctx, cfg Config, doc *Document) error {
	// Filter request
	if cfg.FilterRequest != nil {
		cfg.FilterRequest(c.Request())
	}

	r := NewResolver()
	statusCode, head, err := r.Resolve(c, cfg, doc.HTMLFragment())
	if err != nil {
		return err
	}

	// get final output
	f := doc.HTMLFragment()
	f.AppendHead(cfg.FilterHead(head)...)

	html, err := f.HTML()
	if err != nil {
		return cfg.ErrorHandler(c, err)
	}

	c.Response().SetStatusCode(statusCode)
	c.Response().Header.SetContentType(fiber.MIMETextHTMLCharsetUTF8)
	c.Response().SetBody([]byte(html))

	// Filter response
	if cfg.FilterResponse != nil {
		cfg.FilterResponse(c.Response())
	}

	return nil
}

// readContent opens a named file and read content from it.
func readContent(rf io.ReaderFrom, name string) (n int64, err error) {
	// Read file
	f, err := os.Open(filepath.Clean(name))
	if err != nil {
		return 0, err
	}
	defer func() { err = f.Close() }()
	return rf.ReadFrom(f)
}

// Helper function to set default values.
func configDefault(config ...Config) Config {
	// Init config
	var cfg Config
	if len(config) > 0 {
		cfg = config[0]
	}

	if cfg.ErrorHandler == nil {
		cfg.ErrorHandler = func(c fiber.Ctx, _ error) error {
			return c.Status(fiber.StatusInternalServerError).SendString("cannot create response")
		}
	}

	if cfg.DefaultHost == "" {
		cfg.DefaultHost = "localhost:3000"
	}

	if cfg.FilterHead == nil {
		cfg.FilterHead = func(nodes []*html.Node) []*html.Node {
			return nodes
		}
	}

	if cfg.Client == nil {
		cfg.Client = DefaultClient
	}

	return cfg
}
