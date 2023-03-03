# set version to first arg
version=$1;

GOOS=windows GOARCH=amd64 go build -o bin/bunny_logs.windows.$version.amd64.exe &&
GOOS=linux GOARCH=amd64 go build -o bin/bunny_logs.linux.$version.amd64 &&
GOOS=linux GOARCH=arm64 go build -o bin/bunny_logs.linux.$version.arm64 &&
GOOS=darwin GOARCH=amd64 go build -o bin/bunny_logs.darwin.$version.amd64 &&
GOOS=darwin GOARCH=arm64 go build -o bin/bunny_logs.darwin.$version.arm64