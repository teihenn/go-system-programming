# example: make new ARG=xxx/yyy
.PHONY: new
new:
	mkdir -p ${ARG}
	cd ${ARG} \
		&& go mod init github.com/teihenn/go-system-programming/${ARG} \
		&& touch main.go \
		&& echo "package main \n\nfunc main() {\n\t\n}" > main.go
