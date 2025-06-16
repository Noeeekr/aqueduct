.PHONY: build

build: 
	mkdir -p build/
	cd ./server && go mod tidy 
	cd ./server && go build  -o ../build/aqueduct ./cmd/aqueduct/main.go
	
	cp -r ./web ./build/public
	cp -r ./scripts/install.sh ./build/
	cp -r ./assets/* ./build