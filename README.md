# Go Source Checker

[![Build Status](https://travis-ci.org/daisuzu/gsc.svg?branch=master)](https://travis-ci.org/daisuzu/gsc)
[![Go Report Card](https://goreportcard.com/badge/github.com/daisuzu/gsc)](https://goreportcard.com/report/github.com/daisuzu/gsc)

_gsc_ checkes bits and pieces of problems in Go source code.

## Checks

| Check | Description |
| ----- | ----------- |
| [CtxScope](checker/testdata/ctxscope.go) | Not to use [context.Context](https://golang.org/pkg/context/#Context) outside the scope. |
| [RangePtr](checker/testdata/rangeptr.go) | Not to use pointer to the loop iteration variable. |

## Installation

```sh
go get -u github.com/daisuzu/gsc
```

## Usage

```sh
Usage of gsc:
	gsc [flags] # runs on package in current directory
	gsc [flags] packages
	gsc [flags] directory
	gsc [flags] files... # must be a single package

For more about the flags, see 'gsc -help'.
```
