gotags
======

gotags is a [ctags][]-compatible tag generator for [Go][].

[![Build Status](https://secure.travis-ci.org/jstemmer/gotags.png?branch=master)](http://travis-ci.org/jstemmer/gotags)

Installation
------------

Install or update gotags using the `go get` command:

	go get -u github.com/jstemmer/gotags

Usage
-----

	gotags [options] file(s)

Vim [Tagbar][] configuration
----------------------------

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
