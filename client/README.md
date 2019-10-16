build wasm

```
GOOS=js GOARCH=wasm go build -o main.wasm main.go
```

test

```
live-server
```