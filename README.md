# Go Source Checker

[![Build Status](https://travis-ci.org/daisuzu/gsc.svg?branch=master)](https://travis-ci.org/daisuzu/gsc)
[![Go Report Card](https://goreportcard.com/badge/github.com/daisuzu/gsc)](https://goreportcard.com/report/github.com/daisuzu/gsc)

_gsc_ checkes bits and pieces of problems in Go source code.

## Checks

| Check | Description |
| ----- | ----------- |
| [ctxscope](ctxscope/testdata/src/a/a.go) | Not to use [context.Context](https://golang.org/pkg/context/#Context) outside the scope. |
| [rangeptr](rangeptr/testdata/src/a/a.go) | Not to use pointer to the loop iteration variable. |

## Installation

```sh
go get -u github.com/daisuzu/gsc
```

## Usage

```sh
Usage: gsc [-flag] [package]

Run 'gsc help' for more detail,
 or 'gsc help name' for details and flags of a specific analyzer.
```
