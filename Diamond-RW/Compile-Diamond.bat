@echo off
setlocal enabledelayedexpansion

set GOOS=windows
set GOARCH=386
set CGO_ENABLED=1

garble -literals -seed=random build -ldflags "-H=windowsgui" -ldflags="-w -s" -o "Diamond.exe" .
upx.exe --ultra-brute Diamond.exe