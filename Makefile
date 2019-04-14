genpeg:
	cd fbsparser; peg fbs.peg

.PHONY: example
example: genpeg
	flatc -g -o example  example/example.fbs
	go run main.go  example/example.fbs template/helper.gotmpl example/example

.PHONY: test
test:
	go test ./fbsparser/...
