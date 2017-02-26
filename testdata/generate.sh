#!/bin/sh

set -eu
convert="convert -verbose"

for file in hippos.png; do
  rm -f "${file%.*}-*"
  $convert \
    -gravity center \
    -crop 480x480+0+0 +repage \
    $file png:${file%.*}-crop-square.png
  $convert \
    -gravity northwest \
    -crop 100x100+50+50 +repage \
    $file png:${file%.*}-crop-50,50,150,150.png
  $convert \
    -gravity northwest \
    -crop 160x120+320+240 +repage \
    $file png:${file%.*}-crop-320,240,480,360.png
done
