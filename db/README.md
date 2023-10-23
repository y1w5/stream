# DB

This folder contains a program to generate a SQLite database from Wikipedia
dumps. See https://dumps.wikimedia.org/enwiki/20231020/ 

Run `go run .` to download and generate the database. You will need the Go
compiler and SQLite3 on your machine.


## SQLite performance

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
