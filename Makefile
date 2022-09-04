update-lex.go:
	curl -OL https://raw.githubusercontent.com/golang/go/master/src/text/template/parse/lex.go
	patch -p1 < lex.go.patch
	mv lex.go cmd/lex.go
