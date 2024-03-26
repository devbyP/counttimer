run5: 
	go run . -p custom -s 5

build:
	go build -o ./dist/timer .

install: build
	cp ./dist/timer ~/.local/bin
