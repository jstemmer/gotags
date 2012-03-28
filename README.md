gotags
======

[Exuberant-ctags][] compatible tag generator for [Go][].

Installation
------------
	go get github.com/jstemmer/go-tags
	go install github.com/jstemmer/go-tags

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
			\ 't' : 'type'
		\ },
		\ 'scope2kind' : {
			\ 'type' : 't'
		\ },
		\ 'ctagsbin'  : 'gotags',
		\ 'ctagsargs' : '-sort -silent'
	\ }

[exuberant-ctags]: http://ctags.sourceforge.net
[go]: http://golang.org
[tagbar]: http://majutsushi.github.com/tagbar/
