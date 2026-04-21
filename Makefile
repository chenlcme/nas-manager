.PHONY: build build-docker build-all clean

build:
	cd frontend && npm run build
	rm -rf cmd/server/dist
	cp -r frontend/dist cmd/server/dist
	go build -o nas-manager ./cmd/server

build-docker:
	docker buildx build \
		--platform linux/amd64,linux/arm64 \
		--tag nas-manager:latest \
		--push \
		.

build-all:
	cd frontend && npm run build
	rm -rf cmd/server/dist
	cp -r frontend/dist cmd/server/dist
	GOOS=linux GOARCH=amd64 go build -o nas-manager-linux-amd64 ./cmd/server
	GOOS=linux GOARCH=arm64 go build -o nas-manager-linux-arm64 ./cmd/server
	GOOS=darwin GOARCH=amd64 go build -o nas-manager-darwin-amd64 ./cmd/server
	GOOS=darwin GOARCH=arm64 go build -o nas-manager-darwin-arm64 ./cmd/server

clean:
	rm -f nas-manager*
