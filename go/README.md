# Stream in Go

This folder contains an HTTP server capable of streaming Wikipedia pages in
JSON. Run the following commands to test the application:

1. start the server: `go run .`
2. send HTTP requests: `./all.sh`

## Range function experiment

The latest Go compiler comes with support for iterator:

1. download gotip: `go install golang.org/dl/gotip@latest`
1. download the latest compiler: `gotip download`
1. run the server: `GOEXPERIMENT=rangefunc gotip run -tags=gotip .`

The latest specification improves readability but does not do a lot for speed.

## Benchmark

Run of `all.sh` with Go 1.21:

```
$ go run .
level=INFO msg="listening on 127.0.0.1:8080"
level=INFO msg="incoming request" method=GET url=/v1/pages.list.std status=200 size=578.3MB heap=1.8GB duration=5.278s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.std status=200 size=578.3MB heap=1.8GB duration=5.164s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.std status=200 size=578.3MB heap=1.8GB duration=5.42s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.exp status=200 size=557.4MB heap=727.2MB duration=4.835s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.exp status=200 size=557.4MB heap=777.7MB duration=4.863s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.exp status=200 size=557.4MB heap=758.1MB duration=4.86s
level=INFO msg="incoming request" method=GET url=/v2/pages.list status=200 size=557.4MB heap=680.8MB duration=4.997s
level=INFO msg="incoming request" method=GET url=/v2/pages.list status=200 size=557.4MB heap=721.4MB duration=4.881s
level=INFO msg="incoming request" method=GET url=/v2/pages.list status=200 size=557.4MB heap=746.4MB duration=4.87s
level=INFO msg="incoming request" method=GET url=/v2/pages.stream status=200 size=557.4MB heap=3.8MB duration=5.827s
level=INFO msg="incoming request" method=GET url=/v2/pages.stream status=200 size=557.4MB heap=5.1MB duration=5.913s
level=INFO msg="incoming request" method=GET url=/v2/pages.stream status=200 size=557.4MB heap=3.5MB duration=5.85s
level=INFO msg="incoming request" method=GET url=/v3/pages.list status=200 size=557.4MB heap=835.2MB duration=4.835s
level=INFO msg="incoming request" method=GET url=/v3/pages.list status=200 size=557.4MB heap=626.3MB duration=4.911s
level=INFO msg="incoming request" method=GET url=/v3/pages.list status=200 size=557.4MB heap=692.5MB duration=4.878s
level=INFO msg="incoming request" method=GET url=/v3/pages.stream status=200 size=557.4MB heap=59.4MB duration=4.79s
level=INFO msg="incoming request" method=GET url=/v3/pages.stream status=200 size=557.4MB heap=60.7MB duration=4.794s
level=INFO msg="incoming request" method=GET url=/v3/pages.stream status=200 size=557.4MB heap=60.7MB duration=4.817s
```

Run of `all.sh` with Go tip:

```
GOEXPERIMENT=rangefunc gotip run -tags=gotip .
level=INFO msg="listening on 127.0.0.1:8080"
level=INFO msg="incoming request" method=GET url=/v1/pages.list.std status=200 size=578.3MB heap=2.0GB duration=4.352s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.std status=200 size=578.3MB heap=2.0GB duration=4.076s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.std status=200 size=578.3MB heap=2.0GB duration=4.369s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.exp status=200 size=557.4MB heap=623.0MB duration=3.691s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.exp status=200 size=557.4MB heap=777.0MB duration=3.823s
level=INFO msg="incoming request" method=GET url=/v1/pages.list.exp status=200 size=557.4MB heap=740.3MB duration=3.837s
level=INFO msg="incoming request" method=GET url=/v2/pages.list status=200 size=557.4MB heap=763.3MB duration=3.894s
level=INFO msg="incoming request" method=GET url=/v2/pages.list status=200 size=557.4MB heap=769.9MB duration=3.91s
level=INFO msg="incoming request" method=GET url=/v2/pages.list status=200 size=557.4MB heap=775.4MB duration=3.909s
level=INFO msg="incoming request" method=GET url=/v2/pages.stream status=200 size=557.4MB heap=3.3MB duration=4.878s
level=INFO msg="incoming request" method=GET url=/v2/pages.stream status=200 size=557.4MB heap=6.1MB duration=4.888s
level=INFO msg="incoming request" method=GET url=/v2/pages.stream status=200 size=557.4MB heap=4.8MB duration=4.996s
level=INFO msg="incoming request" method=GET url=/v3/pages.list status=200 size=557.4MB heap=816.0MB duration=3.929s
level=INFO msg="incoming request" method=GET url=/v3/pages.list status=200 size=557.4MB heap=802.2MB duration=3.908s
level=INFO msg="incoming request" method=GET url=/v3/pages.list status=200 size=557.4MB heap=813.9MB duration=3.903s
level=INFO msg="incoming request" method=GET url=/v3/pages.stream status=200 size=557.4MB heap=65.4MB duration=3.923s
level=INFO msg="incoming request" method=GET url=/v3/pages.stream status=200 size=557.4MB heap=66.5MB duration=3.913s
level=INFO msg="incoming request" method=GET url=/v3/pages.stream status=200 size=557.4MB heap=68.7MB duration=3.923s
```
