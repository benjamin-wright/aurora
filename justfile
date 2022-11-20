clean:
    rm -rf dist
    rm -rf bin

build: clean wasm tools

tools:
    mkdir -p bin
    cd tools/img-convert/cmd/normals && go build -o ../../../../bin/normals 

wasm:
    mkdir -p dist
    cd wasm && GOOS=js GOARCH=wasm go build -o ../dist/main.wasm .
    cp "$(go env GOROOT)/misc/wasm/wasm_exec.js" ./dist/wasm_exec.js
    cp ./wasm/index.html ./dist/index.html
    cp -r ./content ./dist/content

start: stop
    docker run --name aurora -p 8080:80 -v $(pwd)/dist:/usr/share/nginx/html:ro -d nginx:alpine

stop:
    docker stop aurora || true
    docker rm aurora || true