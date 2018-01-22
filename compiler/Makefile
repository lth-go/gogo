all : parser.go

parser.go : parser.go.y
	goyacc -o $@ parser.go.y

clean:
	rm -f parser.go y.output
