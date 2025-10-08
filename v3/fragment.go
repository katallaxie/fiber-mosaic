package fragments

import (
	"bytes"
	"context"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

const (
	// HTTPLinkHeader is the HTTP header for LINK entities.
	HTTPLinkHeader = "Link"
)

// HTMLFragment is representation of HTML fragments.
type HTMLFragment struct {
	doc *goquery.Document
	sync.RWMutex
}

// NewHTMLFragment creates a new fragment of HTML.
func NewHTMLFragment(root *html.Node) (*HTMLFragment, error) {
	h := new(HTMLFragment)
	h.doc = goquery.NewDocumentFromNode(root)

	return h, nil
}

// Document get the full document representation
// of the HTML fragment.
func (h *HTMLFragment) Fragment() *goquery.Document {
	return h.doc
}

// Fragments is returning the selection of fragments
// from an HTML page.
func (h *HTMLFragment) Fragments() (map[string]*Fragment, error) {
	h.RLock()
	defer h.RUnlock()

	scripts := h.doc.Find("head script[type=fragment]")
	fragments := h.doc.Find("app-fragment").AddSelection(scripts)

	ff := make(map[string]*Fragment)

	fragments.Each(func(_ int, s *goquery.Selection) {
		f := FromSelection(s)

		if !f.deferred {
			ff[f.ID()] = f
		}
	})

	return ff, nil
}

// HTML creates the HTML output of the created document.
func (h *HTMLFragment) HTML() (string, error) {
	h.RLock()
	defer h.RUnlock()

	html, err := h.doc.Html()
	if err != nil {
		return "", err
	}

	return html, nil
}

// AppendHead ...
func (h *HTMLFragment) AppendHead(ns ...*html.Node) {
	head := h.doc.Find("head")
	head.AppendNodes(ns...)
}

// Fragment is a <fragment> in the <header> or <body>
// of a HTML page.
type Fragment struct {
	deferred bool
	fallback string
	method   string
	primary  bool
	src      string
	timeout  int64

	id  string
	ref string

	statusCode int
	head       []*html.Node

	f *HTMLFragment
	s *goquery.Selection
}

// FromSelection creates a new fragment from a
// fragment selection in the DOM.
func FromSelection(s *goquery.Selection) *Fragment {
	f := new(Fragment)
	f.s = s

	src, _ := s.Attr("src")
	f.src = src

	fallback, _ := s.Attr("fallback")
	f.fallback = fallback

	method, _ := s.Attr("method")
	f.method = method

	timeout, ok := s.Attr("timeout")
	if !ok {
		timeout = "60"
	}
	t, _ := strconv.ParseInt(timeout, 10, 64)
	f.timeout = t

	id, ok := s.Attr("id")
	if !ok {
		id = uuid.New().String()
	}
	f.id = id

	ref, _ := s.Attr("ref")
	f.ref = ref

	deferred, ok := s.Attr("deferred")
	f.deferred = ok && !strings.EqualFold(deferred, "FALSE")

	primary, ok := s.Attr("primary")
	f.primary = ok && !strings.EqualFold(primary, "FALSE")

	f.head = make([]*html.Node, 0)

	return f
}

// Src is the URL for the fragment.
func (f *Fragment) Src() string {
	return f.src
}

// Fallback is the fallback URL for the fragment.
func (f *Fragment) Fallback() string {
	return f.fallback
}

// Timeout is the timeout for fetching the fragment.
func (f *Fragment) Timeout() time.Duration {
	return time.Duration(f.timeout) * time.Second
}

// Method is the HTTP method to use for fetching the fragment.
func (f *Fragment) Method() string {
	return f.method
}

// Element is a pointer to the selected element in the DOM.
func (f *Fragment) Element() *goquery.Selection {
	return f.s
}

// Deferred is deferring the fetching to the browser.
func (f *Fragment) Deferred() bool {
	return f.deferred
}

// Primary denotes a fragment as responsible for setting
// the response code of the entire HTML page.
func (f *Fragment) Primary() bool {
	return f.primary
}

// Links returns the new nodes that go in the head via
// the LINK HTTP header entity.
func (f *Fragment) Links() []*html.Node {
	return f.head
}

// Ref represents the reference to another fragment.
func (f *Fragment) Ref() string {
	return f.ref
}

// ID represents a unique id for the fragment.
func (f *Fragment) ID() string {
	return f.id
}

// HTMLFragment returns embedded fragments of HTML.
func (f *Fragment) HTMLFragment() *HTMLFragment {
	return f.f
}

// Resolve is resolving all needed data, setting headers
// and the status code.
func (f *Fragment) Resolve() ResolverFunc {
	return func(c fiber.Ctx, cfg Config) error {
		err := f.do(c, cfg, f.src)
		if !errors.Is(err, context.DeadlineExceeded) {
			return err
		}

		err = f.do(c, cfg, f.fallback)
		if err != nil {
			return err
		}

		return nil
	}
}

func (f *Fragment) do(c fiber.Ctx, cfg Config, src string) error {
	u, err := url.Parse(src)
	if err != nil {
		return err
	}

	if len(u.Scheme) == 0 {
		u.Scheme = "http"
	}

	if len(u.Host) == 0 {
		u.Host = cfg.DefaultHost
	}

	ctx, cancel := context.WithTimeout(c, f.Timeout())
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return err
	}

	res, err := DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	f.statusCode = res.StatusCode

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	res.Header.Del(fiber.HeaderConnection)

	h := Header(res.Header.Get(HTTPLinkHeader))
	nodes := CreateNodes(h.Links())
	f.head = append(f.head, nodes...)

	root := &html.Node{
		Type:     html.ElementNode,
		DataAtom: atom.Body,
		Data:     "body",
	}

	ns, err := html.ParseFragment(bytes.NewReader(body), root)
	if err != nil {
		return err
	}

	for _, n := range ns {
		root.AppendChild(n)
	}

	doc, err := NewHTMLFragment(root)
	if err != nil {
		return err
	}
	f.f = doc

	return nil
}
