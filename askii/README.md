# askii

tiny ascii art image renderer

## usage

```
./askii [options] <image>
```

## what it does

* converts an image into ascii characters
* maps brightness to a customizable palette of ascii symbols
* optionally uses truecolor to render the original image colors in the terminal
* supports adjustable x/y sampling steps for fine-tuning resolution

## flags

* `-x`: x axis step (default `4`)
* `-y`: y axis step (default `8`)
* `-palette`: ascii palette to use (default `" .:-=+*#%@"`)
* `-gray`: ignore color and render in pure grayscale (default `false`)
