build: main.go
	mkdir bin
	go build -o bin/gonta main.go

run: build
	./bin/gonta
	rm -f bin/gonta

clean:
	rm -rf bin
