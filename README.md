# gotags

gotags is a [ctags][]-compatible tag generator for [Go][].

[![Build Status](https://travis-ci.org/jstemmer/gotags.svg?branch=master)](http://travis-ci.org/jstemmer/gotags)

## Installation

[Go][] version 1.1 or higher is required. Install or update gotags using the
`go get` command:

	go get -u github.com/jstemmer/gotags

## Usage

	gotags [options] file(s)

	-L="": source file names are read from the specified file.
	-R=false: recurse into directories in the file list
	-silent=false: do not produce any output on error
	-sort=true: sort tags
	-tree=false: print syntax tree (debugging)
	-v=false: print version

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

### Screenshot
![vim Tagbar gotags](http://stemmertech.com/images/gotags-1.0.0-screenshot.png)

[ctags]: http://ctags.sourceforge.net
[go]: http://golang.org
[tagbar]: http://majutsushi.github.com/tagbar/
[screenshot]: https://github.com/jstemmer/gotags/gotags-1.0.0-screenshot.png
