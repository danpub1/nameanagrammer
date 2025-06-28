GOOS=js GOARCH=wasm go build -o nameanagrammer.wasm
gzip -9 -k -f -n nameanagrammer.wasm
go build -tags web -o nameanagrammerweb
go build