go-ujson
--------

A pure Go port of [ultrajson][ultrajson] with a [go-simplejson][go-simplejson] like interface.

    $ go test -bench ".*"
    PASS
    BenchmarkUjson	  500000	      4970 ns/op	  20.12 MB/s
    BenchmarkStdLib	  200000	     10323 ns/op	   9.69 MB/s

WARNING: very early stages of a public API

[ultrajson]: https://github.com/esnme/ultrajson
[go-simplejson]: https://github.com/bitly/go-simplejson
