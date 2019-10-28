all:
	go build -o ./build/account main.go
	npm run build --prefix web
