SET GOOS=linux
SET GOARCH=amd64

ECHO GO: Compiling for %GOOS%-%GOARCH%..
go build -ldflags="-w -s" -o build/dutrozkladapi

docker build -t dutrozkladapi .
docker save dutrozkladapi > dutrozkladapi.tar