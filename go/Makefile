export GOPROXY=https://goproxy.io

build: cmd/main.go
	go build -ldflags \
		"-X main.Version=$$(git describe --tags)" \
		"-X main.Repository=$$(git remote -v | grep fetch | awk '{print $2}')" \
		"-X main.Hashcode=$$(git rev-parse HEAD)" \
		"-X main.Datetime=$$(date '+%Y-%m-%d %H:%M:%S')" \
		"-X main.Hostname=$$(hostname):$$(pwd)" \
		"-X main.GoVersion=$$(go version)" \
		cmd/main.go