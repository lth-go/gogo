.PHONY: parser
parser :
	make -C compiler all

.PHONY: test
test :
	go test -v -count=1 ./vm_test.go

.PHONY: build
build :
	go build -o ./bin/gogo main.go
