clean:
	rm -r dist || exit 0
	mkdir -p dist

linux:
	GOOS=linux GOARCH=amd64 go build -o dist/genius-linux-amd64 main.go

windows:
	GOOS=windows GOARCH=amd64 go build -o dist/genius-windows-amd64.exe main.go

macos:
	GOOS=darwin GOARCH=amd64 go build -o dist/genius-darwin-amd64 main.go

macos-silicon:
	GOOS=darwin GOARCH=amd64 go build -o dist/genius-darwin-arm main.go

release: windows linux macos macos-silicon