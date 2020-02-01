#$Env:GOOS="linux"; $Env:GOARCH="arm"; $Env:GOARM="7"; go clean; go build -ldflags "-s -w"

docker buildx build --platform linux/arm/v7 . -t milgradesec/nbot:0.3 -t milgradesec/nbot:latest --push