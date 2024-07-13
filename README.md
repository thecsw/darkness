# I love Darkness 🥬

![darkness](./darkness.png)

<div id='badges' align='center'>

[![Go Report Card](https://goreportcard.com/badge/github.com/thecsw/darkness)](https://goreportcard.com/report/github.com/thecsw/darkness)
[![GoDoc](https://godoc.org/github.com/thecsw/darkness?status.svg)](https://godoc.org/github.com/thecsw/darkness)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

</div>

## Motivation

I have been writing and maintaining [my website](https://sandyuraz.com) for more 
than five years! Asciidoctor was a nice way to bootstrap quick and good-looking
web pages in a matter of minutes. However, my desire for 
[quirky designs](https://sandyuraz.com/blogs/design) and 
[math support](https://sandyuraz.com/blogs/sqrt2irrational), made the issue
of bulding it [annoyingly complicated](https://sandyuraz.com/blogs/web-legacy).
I just want a nice website that **can do anything**.

## Darkness
 
It is my honor and pleasure to introduce you to [Darkness](https://sandyuraz.com/darkness).
The most noble static website generator. To keep the long story short, check out
[my current website](https://sandyuraz.com), which is dutifully built by Darkness.

She supports all org mode (and markdown!) formatting, 
[native flex galleries](https://sandyuraz.com/plastic) 
(with automatic blurry previews generations and lazy loading!), 
[full math support](https://sandyuraz.com/blogs/diffeq),
social embeds ([youtube](https://sandyuraz.com/blogs/best_web),
[spotify](https://sandyuraz.com/blogs/wrapped-2/), etc.),
[drop caps](https://sandyuraz.com/blogs/cameraman) 
([they look fancy](https://support.microsoft.com/en-us/office/insert-a-drop-cap-817fd19f-40fe-4b73-95e8-f3c0f5e01278)),
[automatic code highlighting](https://sandyuraz.com/blogs/mira_reddit), and much more 😈

Also, did I tell you it's **super fast**?

## Performance

```sh
sandyuraz:source λ darkness build
Processed 128 files in 112 ms
farewell
```

*Each* page takes <1ms to process on my 2020 M1 MacBook Pro. With all the IO time included
as well. Heavily optimized with [komi pools](https://github.com/thecsw/komi),
[hunting heap moves](https://hmarr.com/blog/go-allocation-hunting/), and
[heavy profiling](https://sandyuraz.com/blogs/pprof).

You can play with best performance by tuning parallelization parameters with `-workers N`,
and other flags you can use by calling `darkness -help`!

Here is a benchmark with `hyperfine` on the same config as above,

```sh
sandyuraz:source λ hyperfine "darkness build"
Benchmark 1: darkness build
  Time (mean ± σ):     112.3 ms ±   6.6 ms    [User: 413.7 ms, System: 92.3 ms]
  Range (min … max):   103.1 ms … 124.3 ms    26 runs
```

## How to get it

It's simple! If you have go1.22 installed, you can install it through `go` tool with

```sh
go install -v github.com/thecsw/darkness/v3@v3.0.2
```

Or you can also grab pre-built binaries from the 
[releases page](https://github.com/thecsw/darkness/releases).

## Building your Darkness website

Darkness and I provide you with a template website that you can get a copy of 
through the binary! Run the below, which will create a new directory, for example,
here called `ishmael`,

```sh
darkness new ishmael
```

Go and run the website locally with `darkness serve` and explore Darkness in action!

Here is the [web version of ishmael](https://sandyuraz.com/ishmael) to browse around!

Okay, **go, go**! I'll see you later 😘
