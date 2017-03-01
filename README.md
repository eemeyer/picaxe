# Picaxe

Picaxe is an [IIIF 2.1-compliant](http://iiif.io/api/image/2.1/) image server.

# Limitations

* Support for `rotations` other than 0 not yet implemented.
* Support for `quality` other than `default` and `color` not yet implemented.
* Only image processing is implemented. The "Image Information Request" API is not yet implemented.

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

```shell
$ curl http://localhost:7073/api/picaxe/v1/iiif/http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/200,/0/default.png
```

# Features

In addition to IIIF parameters, additional parameters can be specified on the query string. For example:he following features are supported.

```shell
$ curl http://localhost:7073/api/picaxe/v1/iiif/http%3A%2F%2Fi.imgur.com%2FJ1XaOIa.jpg/full/200,/0/default.png?autoOrient=true
```

## Automatic orientation

Images that have EXIF metadata with an orientation can be normalized for older browsers that don't support reorienting them, by passing `autoOrient=true`.

## Border trimming

To enable border trimming, pass the query parameter `trimBorder` set to a value greater or equal to 0, and less than 1.0. This parameter is the "fuzz factor".

The edge of the image is considered a trimmable border iff it is contiguous with respect to color distance. A color is contiguious iff the distance to the adjacent pixel's color is less than or equal to the fuzz factor. (With a fuzz factor of 0.0, all colors are distinct.) Furthermore, the border must extend around the entire rectangular edge of the image. The algorithm trims the outer edge concentrically until a non-consecutive edge is found.

# License

BSD. See `LICENSE` file.
