go-tags
===============

Go tags generator

Installation
------------

	go get github.com/jstemmer/go-tags

	go install github.com/jstemmer/go-tags

Vim Tagbar configuration
------------------------
	let g:tagbar_type_go = {
		\ 'ctagstype' : 'go',
		\ 'kinds'     : [
			\ 'p:package',
			\ 'i:imports'
		\ ],
		\ 'ctagsbin'  : 'gotags',
		\ 'ctagsargs' : ''
	\ }
