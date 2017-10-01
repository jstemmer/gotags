# gotags

gotags is a [ctags][]-compatible tag generator for [Go][].

[![Build Status][travis-badge]][travis-link]
[![Report Card][report-badge]][report-link]

## Installation

[Go][] version 1.1 or higher is required. Install or update gotags using the
`go get` command:

	go get -u github.com/jstemmer/gotags

Or using package manager `brew` on OS X

	brew install gotags

## Usage

	gotags [options] file(s)

	-L="": source file names are read from the specified file. If file is "-", input is read from standard in.
	-R=false: recurse into directories in the file list.
	-f="": write output to specified file. If file is "-", output is written to standard out.
	-silent=false: do not produce any output on error.
	-sort=true: sort tags.
	-tag-relative=false: file paths should be relative to the directory containing the tag file.
	-v=false: print version.

## Vim [Tagbar][] configuration

Put the following configuration in your vimrc:

	let g:tagbar_type_go = {
		\ 'ctagstype' : 'go',
		\ 'kinds'     : [
			\ 'p:package',
			\ 'i:imports:1',
			\ 'c:constants',
			\ 'v:variables',
			\ 't:types',
			\ 'n:interfaces',
			\ 'w:fields',
			\ 'e:embedded',
			\ 'm:methods',
			\ 'r:constructor',
			\ 'f:functions'
		\ ],
		\ 'sro' : '.',
		\ 'kind2scope' : {
			\ 't' : 'ctype',
			\ 'n' : 'ntype'
		\ },
		\ 'scope2kind' : {
			\ 'ctype' : 't',
			\ 'ntype' : 'n'
		\ },
		\ 'ctagsbin'  : 'gotags',
		\ 'ctagsargs' : '-sort -silent'
	\ }

### Vim+Tagbar Screenshot
![vim Tagbar gotags](https://stemmertech.com/images/gotags-1.0.0-screenshot.png)

## gotags with Emacs

Gotags doesn't have support for generating etags yet, but
[gotags-el](https://github.com/craig-ludington/gotags-el) allows you to use
gotags directly in Emacs.

[ctags]: http://ctags.sourceforge.net
[go]: https://golang.org
[tagbar]: https://majutsushi.github.com/tagbar/
[screenshot]: https://github.com/jstemmer/gotags/gotags-1.0.0-screenshot.png
[travis-badge]: https://travis-ci.org/jstemmer/gotags.svg?branch=master
[travis-link]: https://travis-ci.org/jstemmer/gotags
[report-badge]: https://goreportcard.com/badge/github.com/jstemmer/gotags
[report-link]: https://goreportcard.com/report/github.com/jstemmer/gotags
