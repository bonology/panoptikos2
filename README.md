# Panoptikos

## Introduction

Panoptikos is an image viewer for Reddit. This repo contains a new, experimental, version I'm working on. The server-side code is written in Go, the client-side heavily relies on JavaScript.

Once the biggest problems are fixed and the most requested features are added, this version will replace the current version of Panoptikos running at http://panoptikos.com/.

## What’s different?

* You can select any subreddits you want (right now only by manipulating the URL in your browser address bar).
* You can save the URL to save your selection of subreddits.
* By default, you see content from Reddit's ”Hot” section instead of from the ”New” section. This protects us largely from trolls that post NSFL content.

## Known Issues

* Uneven column bottoms because images are not distributed by column height.
* Possible high memory usage and even browser crashes on 32-bit systems because images outside the visible screen area are not removed.
* Jerky scrolling because high-resolution images from sources other than Imgur are often used for the small previews.
* Resizing the browser window blanks the page.
* Links to the Reddit thread are not perceived as links but rather normal text.
* There is no endless scrolling, you must click "Load more" all the time.
* There is no convenient way to select subreddits.
* You can't switch between Reddit's "Hot", "New", "Controversial", "Top" sections. (Current default is "Hot".)

## Compiling Panoptikos

[Install the Go tools](http://golang.org/doc/install) if you haven't done so already. Then, clone the repository, go to the cloned directory and build the executable:

	$ git clone https://github.com/ChristianSiegert/panoptikos.git
	$ cd ./panoptikos
	$ go build

The executable is named after the directory, in this case "panoptikos".

## Running Panoptikos

Simply execute the compiled file:

	$ ./panoptikos

## Command-line arguments

Panoptikos supports these command-line arguments:

* **--port** HTTP port the web server listens to. Default is "8080".

## Example

To start the web server and make it listen to port 80:

	$ ./panoptikos --port=80