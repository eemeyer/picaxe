package main

import (
	"bytes"
	"io"

	"github.com/rwcarlsen/goexif/exif"
)

type Metadata struct {
	Exif *exif.Exif
}

func NewMetadataFromReader(r io.Reader) Metadata {
	x, _ := exif.Decode(r)
	return Metadata{Exif: x}
}

func NewMetadataFromBytes(b []byte) Metadata {
	return NewMetadataFromReader(bytes.NewReader(b))
}
