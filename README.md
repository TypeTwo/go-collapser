Collapser
===

[![GoDoc](https://godoc.org/github.com/cristiangreco/go-collapser?status.svg)](https://godoc.org/github.com/cristiangreco/go-collapser)
[![Build Status](https://travis-ci.org/cristiangreco/go-collapser.svg?branch=master)](https://travis-ci.org/cristiangreco/go-collapser)
[![Go Report Card](https://goreportcard.com/badge/github.com/cristiangreco/go-collapser)](https://goreportcard.com/report/github.com/cristiangreco/go-collapser)

Package `go-collapser` implements a function call deduplication utility.

---

## Install

```sh
go get github.com/cristiangreco/go-collapser
```

## Design goals

- **No external dependencies** - only stdlib packages 
- **Easy to use** - simple, tested and documented API
- **Small code base** - do one thing and do it well
- **Composable** - easy to plug into different contexts (http handlers, db access, remote apis, ...)

## License

The Collapser source files are distributed under the BSD-style license found in the [LICENSE](https://github.com/cristiangreco/go-collapser/blob/master/LICENSE) file.
