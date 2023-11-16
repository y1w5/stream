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

## Summary performace

Adding the `Summarize` function makes the program 6 times slower:

```
$ time go run .
Loading Wikipedia dataset...
Setting up SQLite database...
Loading dataset into SQLite...
...
Completed, 27379 pages created.

real	1m26,923s
user	1m28,145s
sys	0m1,136s
```

Performance can certainly be improved, see:

- https://dave.cheney.net/paste/gophercon-sg-2023.html
- https://research.swtch.com/pcdata
