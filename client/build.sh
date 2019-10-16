GOOS=js GOARCH=wasm go build -o main.wasm main.go
mv main.wasm ../server/public