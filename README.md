# go-gsrv

[![GoDoc](https://godoc.org/github.com/weathersource/go-gsrv?status.svg)](https://godoc.org/github.com/weathersource/go-gsrv)
[![Go Report Card](https://goreportcard.com/badge/github.com/weathersource/go-gsrv)](https://goreportcard.com/report/github.com/weathersource/go-gsrv)
[![Build Status](https://travis-ci.org/weathersource/go-gsrv.svg)](https://travis-ci.org/weathersource/go-gsrv)
[![Codevov](https://codecov.io/gh/weathersource/go-gsrv/branch/master/graphs/badge.svg)](https://codecov.io/gh/weathersource/go-gsrv)

Package gsrv creates a test gRPC server that is useful when mocking a gRPC service.

A complete [example of mocking a gRPC service](https://github.com/weathersource/go-gsrv/tree/master/examples/foo) that leverges this package is available.

This package is based on Google's testutil.Server. Original code: https://github.com/GoogleCloudPlatform/google-cloud-go/blob/master/internal/testutil/server.go
