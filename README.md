# Picaxe

Picaxe is an [IIIF 2.1-compliant](http://iiif.io/api/image/2.1/) image server.

# Limitations

* Support for `rotations` other than 0 not yet implemented.
* Support for `quality` other than `default` and `color` not yet implemented.

# Requirements

* Go 1.7 or later.
* GNU Make to build.

# Running

Build with:

```shell
$ make
```

Then run:

```shell
$ ./build/picaxe -l localhost
```

This will start it on the default port 7073. You can now try a URL such as:

```
$ curl http://localhost:7073/api/picaxe/v1/iiif/http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/200,/0/default.png
```

# License

BSD. See `LICENSE` file.
