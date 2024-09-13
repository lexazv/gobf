SRC_PATH=$(PWD)

build: main.go
	go build -o gobf ${SRC_PATH}/main.go
