SET GOOS=windows
SET GOARCH=amd64

ECHO GO: Compiling for %GOOS%-%GOARCH%..
go build -ldflags="-w -s" -o build/dutrozklad_api.exe

REM docker build -t dutrozkladapi .
REM docker save dutrozkladapi > dutrozkladapi.tar