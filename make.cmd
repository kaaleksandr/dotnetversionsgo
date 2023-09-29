@echo off
chcp 65001

REM 
REM WINDOWS (x86)
REM 

set GOOS=windows
set GOARCH=386
echo Building (Windows_x86) dotnetversionsgo.exe
go build -o bin/x86/dotnetversionsgo.exe -ldflags "-s -w"

set GOARCH=amd64
echo Building (Windows_x64) dotnetversionsgo.exe
go build -o bin/x64/dotnetversionsgo.exe -ldflags "-s -w"

set GOARCH=arm
echo Building (Windows_arm) dotnetversionsgo.exe
go build -o bin/arm/dotnetversionsgo.exe -ldflags "-s -w"

set GOARCH=arm64
echo Building (Windows_arm64) dotnetversionsgo.exe
go build -o bin/arm64/dotnetversionsgo.exe -ldflags "-s -w"
