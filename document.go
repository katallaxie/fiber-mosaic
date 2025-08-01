package fragments

import (
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/net/html"
)

// Document ...
type Document struct {
	doc        *goquery.Document
	html       *HTMLFragment
	statusCode int

	sync.RWMutex
}

// NewDocument ...
func NewDocument(root *html.Node) (*Document, error) {
	d := new(Document)
	// set the default status code
	d.statusCode = fiber.StatusOK

	html, err := NewHTMLFragment(root)
	if err != nil {
		return nil, err
	}
	d.html = html

	return d, nil
}

// HTML is returning the final HTML output.
func (d *Document) HTML() (string, error) {
	d.RLock()
	defer d.RUnlock()

	html, err := d.html.HTML()
	if err != nil {
		return "", err
	}

	return html, nil
}

// Fragments is returning the selection of fragments
// from an HTML page.
func (d *Document) Fragments() ([]*Fragment, error) {
	d.RLock()
	defer d.RUnlock()

	scripts := d.doc.Find("head script[type=fragment]")
	fragments := d.doc.Find("fragment").AddSelection(scripts)

	ff := make([]*Fragment, 0, fragments.Length())

	fragments.Each(func(_ int, s *goquery.Selection) {
		f := FromSelection(s)

		if !f.deferred {
			ff = append(ff, f)
		}
	})

	return ff, nil
}

// HTMLFragment returns embedded fragments of HTML.
func (d *Document) HTMLFragment() *HTMLFragment {
	d.RLock()
	defer d.RUnlock()

	return d.html
}

// SetStatusCode is setting the HTTP status code for the document.
func (d *Document) SetStatusCode(status int) {
	d.Lock()
	defer d.Unlock() // could do this atomic

	d.statusCode = status
}

// StatusCode is getting the HTTP status code for the document.
func (d *Document) StatusCode() int {
	d.RLock()
	defer d.RUnlock()

	return d.statusCode
}
