all:
	go build -buildmode=c-shared -o dist/signalfx.so

clean:
	rm -rf dist