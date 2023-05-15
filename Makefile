# if you are using WSL, build the exe for windows
# and execute exe directly as
# native Ubuntu clipboard performance of in WSL is very poor
exe:
	GOOS=windows go build .

release:
	goreleaser build --snapshot
