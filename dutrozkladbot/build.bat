SET GOOS=windows
SET GOARCH=amd64

ECHO GO: Compiling for %GOOS%-%GOARCH%..
go build -ldflags="-w -s" -o build/dutrozklad_bot.exe

REM docker build -t dutrozkladbot .
REM docker save dutrozkladbot > dutrozkladbot.tar