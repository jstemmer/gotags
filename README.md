gotags
======

[Exuberant-ctags][] compatible tag generator for [Go][].

Installation
------------

	go get github.com/jstemmer/go-tags

Usage
-----

	gotags [options] file(s)

Vim [Tagbar][] configuration
----------------------------------------------------------------
	let g:tagbar_type_go = {
		\ 'ctagstype' : 'go',
		\ 'kinds'     : [
			\ 'p:package',
			\ 'i:imports:1',
			\ 'c:constants',
			\ 'v:variables',
			\ 't:types',
			\ 'w:fields',
			\ 'f:funcs'
		\ ],
		\ 'sro' : '.',
		\ 'kind2scope' : {
			\ 't' : 'ctype'
		\ },
		\ 'scope2kind' : {
			\ 'ctype' : 't'
		\ },
		\ 'ctagsbin'  : 'gotags',
		\ 'ctagsargs' : '-sort -silent'
	\ }

[exuberant-ctags]: http://ctags.sourceforge.net
[go]: http://golang.org
[tagbar]: http://majutsushi.github.com/tagbar/
