@echo off
chcp 65001

REM 
REM WINDOWS (x86)
REM 

set GOOS=windows
set GOARCH=386

echo Building (Windows_x86) dotnetversionsgo.exe

go build -o bin/dotnetversionsgo.exe -ldflags "-s -w"
