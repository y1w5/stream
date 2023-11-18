# DB

This folder contains a program to generate a SQLite database from Wikipedia
dumps. See https://dumps.wikimedia.org/enwiki/20231020/ 

Run `go run .` to download and generate the database. You will need the Go
compiler and SQLite3 on your machine.


## SQLite performance

Using transactions speed-up database creation by a factor 30 (7m30s to 13s).

Without transaction:

```
$ time go run .
Loading Wikipedia dataset...
Setting up SQLite database...
Loading dataset into SQLite...
...
Completed, 27379 pages created.

real	7m18,788s
user	1m0,634s
sys	0m23,778s
```

With transaction:

```
$ time go run .
Loading Wikipedia dataset...
Setting up SQLite database...
Loading dataset into SQLite...
...
Completed, 27379 pages created.

real	0m13,391s
user	0m12,970s
sys	0m1,958s
```

## XML Performance

I wanted to benchmark the `Summarize` function and I discovered that the `xml`
package is quite slow. Reading a page from the dataset takes 0.4ms and allocates
80ko 😱.

I tried to speed up the decoding using recommendations from Stackoverflow, but
the code is overly complicated and the gains are limited. I also added a function
to benchmark the decoder while streaming the bzip2 archive.

```
goos: linux
goarch: amd64
pkg: github.com/y1w5/stream/db
cpu: 12th Gen Intel(R) Core(TM) i7-1260P
BenchmarkDecoder-16              	   5287	   405179 ns/op	  84046 B/op	    158 allocs/op
BenchmarkDecoderV2-16            	   5511	   377566 ns/op	  42457 B/op	    106 allocs/op
BenchmarkDecoder_streaming-16    	   1416	  1620417 ns/op	  86662 B/op	    160 allocs/op
BenchmarkSummarize-16            	   5148	   408327 ns/op	  88073 B/op	    166 allocs/op
PASS
ok  	github.com/y1w5/stream/db	39.271s
```

- Stackoverflow: https://stackoverflow.com/a/61858457
