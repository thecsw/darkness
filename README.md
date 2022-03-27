# I love Darkness 🥬

Also [posted here](https://sandyuraz.com/darkness/)

![darkness](./darkness.png)

<div id='badges' align='center'>

[![Go Report Card](https://goreportcard.com/badge/github.com/thecsw/darkness)](https://goreportcard.com/report/github.com/thecsw/darkness)
[![GoDoc](https://godoc.org/github.com/thecsw/darkness?status.svg)](https://godoc.org/github.com/thecsw/darkness)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

</div>


## This is no life

It doesn\'t feel right anymore. I feel that every time I write something
for my website, I have to dance to the tune of `pandoc` and
`asciidoctor`. Maybe the Ruby and Haskell Gods will take pity
on me and my stuff will actually render as I want them to. Oh... they
didn\'t? Time to write a new `sed` script to fix that.

This is no life. Instead of fixing an issue at its root, I\'ve been
writing
[fixes](https://github.com/thecsw/thecsw.github.io/blob/legacy-source/sed/html.sed),
[patches](https://github.com/thecsw/thecsw.github.io/blob/legacy-source/sed/adoc.sed),
and
[this](https://github.com/thecsw/thecsw.github.io/blob/legacy-source/Makefile).
It was fun in the beginning! Mastering sed and awk, the forbidden tools
echoed from the past. However, as time went on and I wanted to add more
content and get more freedom, [this
pipeline](https://sandyuraz.com/blogs/web-legacy/) felt limiting.

## A hero emerges

It is my honor and pleasure to introduce you to
[Darkness](https://github.com/thecsw/Darkness)! The most noble static
website generator. The website you are on right now is built by her,
with all the images, styling, embeds, and more! The big advantage of
Darkness is that she is general enough (with a config file), so she can
build any website for what their heart desires.

Darkness currently takes a subset of orgmode files and creates an html
file for them. Orgmode is a markup language used by emacs, though it\'s
very similar to markdown. There is a getting started portion further! It
will get you everything you need to get comfortable with the process.
Also, darkness is more agnostic than just org-\>html. One just has to
write a parser/exporter for a markup language, enable it, and now we
have multiple supported types!

## Plethora of features

Darkness supports a vast array of features and cool tricks, such as

-   Project configuration with a [simple toml
    file](https://github.com/thecsw/darkness/blob/master/ishmael/darkness.toml)
-   Smart headings (laid out and shifted as needed)
-   Paragraphs with full [orgmode-style
    formatting](https://orgmode.org/worg/dev/org-syntax.html)
-   Lists with any text an breaking lines
-   Good placement of images and its captions
-   Embed support for: youtube, spotify, video, audio, etc.
-   Raw HTML placement and export within the page
-   Source code blocks with language selection
-   Full math render support with [KaTeX](https://katex.org) (only
    enabled on pages with math)
-   Smart Holoscene time annotations and layovers (just hover over H.E.
    times)
-   Rich meta tags (OG/twitter-friendly) with previews enabled on every
    page
-   Creating new darkness websites from a template (use
    `darkness new NAME`)
-   Project cleanup and maintenance tools included with the binary
-   Open-source code, PRs and contributions are welcome!

## Not only noble, but also super fast

Darkness is also lightning fast! With the old pipeline, it would take
around 8-10 seconds to build my website. Darkness **only takes \~100ms**
to do the same, with IO and memory syscall times included. This is why
[profiling should be your friend.](https://sandyuraz.com/blogs/pprof/)
She uses some cool parallelization tricks, so we spend around 1ms on
each page. Pretty cool if you ask me. Fastest (though admittedly larger)
static website generators took a bit more time in my experience.

``` {.bash org-language="sh"}
sandyuraz:source* λ darkness build
Looking for files... found 83 in 23 ms
Building and flushing... done in 65 ms
farewell
```

## Getting started

Want to try out darkness but don\'t know orgmode? You are at the very
start of something beautiful, my friend! If you have go 1.18+ installed,
run the command below to get the latest version of her (you can also
grab a suitable binary from the
[releases](https://github.com/thecsw/darkness/releases) page)

``` {.bash org-language="sh"}
go install -v github.com/thecsw/darkness@latest
```

Next, go to the place, where you would like to hous your next website.
Call her with a `new` command followed by any name that feeds
your soul (`ishmael` is an example here)

``` {.bash org-language="sh"}
darkness new ishmael
```

Voilà! You have a new darkness project initiated. Change into that new
directory with `cd ishmael` and type in
`darkness build`. You\'ll find an `index.html` in
the root, go open it! After that, look into `darkness.toml`,
which shows you the number of things you can play with to change the
looks and feel.

Here is [Ishmael\'s website](https://sandyuraz.com/ishmael) that shows
the formatting in action with many other cool things! It\'s the template
website that we just made above. I do hope you find it fun here

## Formatting

Here is some orgmode and darkness formatting guide! The actual
formatting and final layout is [here](https://sandyuraz.com/ishmael)!

### Sections

Every darkness input should start with the page\'s title. Page\'s title
will look in the form of `* Page's title` and every section
on the website will start with two or more astericks to show the heading
level. So, `** Section` and `***
Subsection` will be put as a child of the former.

### Formatting

You can separate text into paragraphs by adding a new line between them.
Think of an empty line as a delimeter. You should always add it between
different types of blocks you add to your page.

Simple formatting is supported as well, you can do some **bold** text,
maybe even *italics*, and `verbatim`. For
`*bold*`, surround your text with astericks, for
`/italics/`, surround with forward slashes, and verbatim with
equal signs, `=verbatim=`

### Links and embeds

Links are in the form of `[­[link][text]]`. If your link is
in a text, then it will show up as
[such](https://en.wikipedia.org/wiki/Ishmael_(Moby-Dick)). If the link
is on a separate line, darkness will decide if it can be inserted as an
embed (image, youtube, spotify song/playlist, video, audio).

### Images

See for example, below is a link on its line with an image path

![*CUT*, August 2009 issue. Art by Tadashi Hiramatsu](https://sandyuraz.com/ishmael/evangelion.webp)

### Songs

Depending on the link type, darkness will intelligently stub in a
preview if it\'s a standalone link (not inlined within text). One more
example

[Last Surprise](https://open.spotify.com/track/4cPnNnTMkJ6soUOUzEtmcp?si=ba1730fdb66642b9)

### Lists

Lists are created by starting a line with a hyphen followed by an item
description, you would write something like

    - This is my first item
    - This second item is going to be so long that I would
    have to break it down into two line
    - Third item follows swiftly

It will render as follows

-   This is my first item
-   This second item is going to be so long that I would have to break
    it down into two line
-   Third item follows swiftly

### Source code blocks

Source code blocks follow
[orgmode](https://orgmode.org/manual/Working-with-Source-Code.html)\'s
conventions. You would wrap your source code with
`#+begin_src LANG` and `#+end_src`, where
`LANG` is the language of the source code block. You can
leave `LANG` empty as well.

``` org
#+begin_src c
main( ) {
        printf("hello, world");
}
#+end_src
```

Will render as (hover over the block to see the language)

``` c
main( ) {
        printf("hello, world");
}
```

### HTML injection

Whatever darkness provides can still be limiting if you want to insert
some of your own material or embeds that are not supported. Do you give
up? Hopefully not. Anything surrounded by
`#+begin_export html` and `#+end_export` will be
inserted literally into the page.

```{=html}
<script src="https://gist.github.com/thecsw/c80f83c0d52c0a476e86fc9a6a980517.js"></script>
```
This is the embed source for the above

``` org
#+begin_export html
<script src="https://gist.github.com/thecsw/c80f83c0d52c0a476e86fc9a6a980517.js"></script>
#+end_export
```

### Attention blocks

You may need to grab reader\'s attention even more or make them aware of
some sharp edges in whatever you\'re writing about. Start a paragraph
with `NOTE:`, `IMPORTANT:`, `CAUTION:`,
`TIP:`, or `WARNING:`, and you will get an
attention-grabber

TIP: This is kinda useful for technical posts when you mention
exceptions or so

### Footnotes

Another cool thing darkness can do for you is keeping track of your
footnotes. We follow orgmode\'s conventions as well. Anywhere in the
text, if you have a string in the form of
`[fn­:: blablabla]`, that `blablabla` will go and
become your footnote. [^1]

## Niche features

See this cool trick of cleaning up the project with
`megumin`!

[My name is Megumin, the number one mage of Axel!](https://sandyuraz.com/darkness/megumin.mp4)

## Why \"Darkness\"?

Her name is based one one of the characters I love from [KonoSuba](https://en.wikipedia.org/wiki/KonoSuba) 

![Dustiness Ford Lalatina](https://sandyuraz.com/darkness/darkness.webp)

*A knight must never run away, no matter how mighty the enemy.* --
Darkness

[^1]: *Formatting* **also** `works` in footnotes
