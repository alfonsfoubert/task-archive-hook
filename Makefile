HOOKS_DIR := ~/.task/hooks
BINARY_NAME := on-modify.task-archive

build:
	go build -o=${BINARY_NAME} main.go

clean:
	rm -f ${HOOKS_DIR}/${BINARY_NAME}  

install: clean build
	mv ${BINARY_NAME} ${HOOKS_DIR}
