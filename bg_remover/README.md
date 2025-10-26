# bg_remover

tiny background remover

## usage

```
./bg_remover -in input.png -out output.png -t 60 -s 8 -soft
```

## what it does

* samples colors from the four corners of the image to estimate a uniform background
* compares each pixel to that average color
* makes pixels near the background color transparent
* runs each image slice concurrently cuz we don't want unefficient and slow software

## flags

* `-in`: input png file (default `in.png`)
* `-out`: output png file (default `out.png`)
* `-t`: color distance threshold (default 60)
* `-s`: corner sample size (default 8)
* `-soft`: enables soft alpha transition (default true)
