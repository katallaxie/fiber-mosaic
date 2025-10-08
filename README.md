# 🖼️ Mosaic

[![Test & Build](https://github.com/katallaxie/fiber-mosaic/actions/workflows/main.yml/badge.svg)](https://github.com/katallaxie/fiber-mosaic/actions/workflows/main.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/katallaxie/fiber-mosaic)](https://goreportcard.com/report/github.com/katallaxie/fiber-mosaic)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Taylor Swift](https://img.shields.io/badge/secured%20by-taylor%20swift-brightgreen.svg)](https://twitter.com/SwiftOnSecurity)

Mosaic middleware for [Fiber](https://github.com/gofiber/fiber) enables building microservices for the frontend.

## 🛸 Installation

```bash
go get github.com/katallaxie/fiber-mosaic/v3
```

## Usage

A `<app-fragment>` symbolizes a part of a template that can be served by a singular microservice. Thus, making a fragment the contract between different services and teams within a large engineering organization. The middleware concurrently fetches those parts from the service and replaces it in the template. It supports `GET` and `POST` [HTTP methods](https://developer.mozilla.org/de/docs/Web/HTTP/Methods) to fetcht the content. Related resources like CSS or JavaScript are injected via the [HTTP `LINK` entity header field](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Link). A `<app-fragment>` can occure in the [`body` element](https://developer.mozilla.org/de/docs/Web/HTML/Element/body) or the [`header` element](https://developer.mozilla.org/de/docs/Web/HTML/Element/header). See [Example](#example) to learn more about using fragments.

`app-fragment` is used to adhere to the [HTML spec](https://html.spec.whatwg.org/#valid-custom-element-name), the tag name must contain a dash ('-').

[Tailor](https://github.com/zalando/tailor) by Zalando is prior art for this middleware.
[Fragements](https://github.com/github/fiber-fragments) by GitHub is prior art for this middleware created by @katallaxie.

## Fragement(s)

A `app-fragment` will be hybrid-polymorphic (if this is a thing). On the server it is parsed and evaluate by the middleware. 🦄 In the browser it will be a web component that received data from the middleware (**this is still work in progress ⚠️**).

### Server

* `src` The source to fetch for replacement in the DOM
* `method` can be of `GET` (default) or `POST`.
* `primary` denotes a fragment that sets the response code of the page
* `id` is an optional unique identifier (optional)
* `ref`is an optional forward reference to an `id` (optional)
* `timeout` timeout of a fragement to receive in milliseconds (default is `300`)
* `deferred` is deferring the fetch to the browser
* `fallback` is the fallback source in case of timeout/error on the current fragment


## Example

Import the middleware package this is part of the Fiber web framework.

```go
package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"

	"github.com/katallaxie/fiber-mosaic"
)
```

After you initiate your Fiber app, you can plugin in the fragments middleware. The middleware draws the templates for the fragments to load from the template engine. Thus it supports using all [template](https://github.com/gofiber/template) engines supported by the Fiber team.

```go
// Create a new engine
engine := html.New("./views", ".html")

// Pass the engine to the Views
app := fiber.New(fiber.Config{
	Views: engine,
})

// Associates the route with a specific template with fragments to render
app.Get("/index", fragments.Template(fragments.Config{}, "index", fiber.Map{}, "layouts/main"))

// this would listen to port 8080
app.Listen(":8080")
```

```html
<html>
<head>
    <script type="fragment" src="assets"></script>
</head>
<body>
    <h1>Example</h1>
    <app-fragment src="fragment1.html"></app-fragment>
</body>
</html>
```

The `example` folder contains many examples. You can learn how to use a forward reference of content for fetching outer and inner content and replace either one.

## Benchmark(s)

This is run on a MacBook Pro 16 inch locally. It is the `example` run.

* Parsing a local template with extrapolation with the fragments
* Parsing the fragments
* Doing fragments
* Inlining results and adding `Link` header resources to the output

```bash
echo "GET http://127.0.0.1:8080/index" | vegeta attack -duration=5s -rate 2000 | tee results.bin | vegeta report 
Requests      [total, rate, throughput]         10000, 2000.22, 1999.86
Duration      [total, attack, wait]             5s, 4.999s, 898.542µs
Latencies     [min, mean, 50, 90, 95, 99, max]  278.625µs, 1.549ms, 805.833µs, 1.591ms, 7.847ms, 16.35ms, 23.643ms
Bytes In      [total, mean]                     10130000, 1013.00
Bytes Out     [total, mean]                     0, 0.00
Success       [ratio]                           100.00%
Status Codes  [code:count]                      200:10000  
Error Set:
```

## License

[MIT](/LICENSE)
