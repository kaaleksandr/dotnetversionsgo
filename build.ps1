Set-Variable -Name GOOS -Value windows
Set-Variable -Name GOARCH -Value 386

#
go build -o bin\${GOOS}\${GOARCH}\ -ldflags '-s' .


Set-Variable -Name GOARCH -Value amd64

#
go build -o bin\${GOOS}\${GOARCH}\ -ldflags '-s' .

Remove-Variable GOOS,GOARCH
