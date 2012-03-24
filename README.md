gotags
======

Go tags generator

Installation
------------
	go get github.com/jstemmer/go-tags
	go install github.com/jstemmer/go-tags

Vim [Tagbar](http://majutsushi.github.com/tagbar/) configuration
------------------------
	let g:tagbar_type_go = {
		\ 'ctagstype' : 'go',
		\ 'kinds'     : [
			\ 'p:package',
			\ 'i:imports:1',
			\ 's:structs',
			\ 'f:functions'
		\ ],
		\ 'ctagsbin'  : 'gotags',
		\ 'ctagsargs' : ''
	\ }
