.PHONY: release test

FILES := $(shell find ./tests/ -name "*.bf" -exec basename -s .bf {} \;)
JOBS := $(addprefix job,${FILES})

test: ${JOBS} ; @echo "[$@] finished!"

${JOBS}: job%: ; go run main.go tests/$*.bf

release:
	GOOS=windows GOARCH=amd64 go build -o ./bin/go-bff_windows_amd64
	GOOS=linux GOARCH=amd64 go build -o ./bin/go-bff_linux_amd64
	GOOS=darwin GOARCH=amd64 go build -o ./bin/go-bff_darwin_amd64
