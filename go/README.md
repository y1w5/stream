# Stream in Go

This folder contains an HTTP server capable of streaming Wikipedia pages in
JSON.

## Range function experiment

The latest Go compiler comes with support for iterator:

1. download gotip: `go install golang.org/dl/gotip@latest`
1. download the latest compiler: `gotip download`
1. run the server: `GOEXPERIMENT=rangefunc gotip run -tags=gotip .`

The latest specification improves readability but does not do a lot for speed.
