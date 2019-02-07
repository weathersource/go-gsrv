# go-gsrv

[![Build Status](https://travis-ci.com/weathersource/go-gsrv.svg?branch=master)](https://travis-ci.com/weathersource/go-gsrv)
[![Codevov](https://img.shields.io/codecov/c/github/weathersource/go-gsrv.svg)](https://codecov.io/gh/weathersource/go-gsrv)
[![Go Report Card](https://goreportcard.com/badge/github.com/weathersource/go-gsrv)](https://goreportcard.com/report/github.com/weathersource/go-gsrv)
[![GoDoc](https://img.shields.io/badge/godoc-ref-blue.svg)](https://godoc.org/github.com/weathersource/go-gsrv)

Package gsrv creates a test gRPC server that is useful when mocking a gRPC service.

A complete [example of mocking a gRPC service](https://github.com/weathersource/go-gsrv/tree/master/examples/foo) that leverges this package is available.

This package is based on Google's testutil.Server. Original code: https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/internal/testutil/server.go
