GOOS=js GOARCH=wasm go build -o nameanagrammer.wasm
gzip -9 -k -f -n nameanagrammer.wasm
go build -tags web -o nameanagrammerweb
go build
base64 nameanagrammer.wasm.gz | gawk '{print $0 "\\"}' >nameanagrammer.wasm.gz.b64
gawk -f embed.awk nameanagrammer.html.src >nameanagrammer.html
rm nameanagrammer.wasm.gz.b64