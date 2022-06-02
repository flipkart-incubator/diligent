# Configure targets
.PHONY: compile clean test start-minion stop-minion start docker
default: compile

BUILD_DIR = ./build
MINION_APP = ${BUILD_DIR}/minion
SHELL_APP = ${BUILD_DIR}/shell

PID_FILE := /tmp/.minion.pid

compile:
	@echo "Compiling..."
	@mkdir -p ${BUILD_DIR}
	@go build -o ${MINION_APP} github.com/flipkart-incubator/diligent/apps/minion
	@go build -o ${SHELL_APP} github.com/flipkart-incubator/diligent/apps/shell
	@echo "Done"

clean:
	@echo "Cleaning..."
	@-rm -f ${MINION_APP} ${SHELL_APP}
	@-rmdir ${BUILD_DIR} 2> /dev/null
	@echo "Done"

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done"

start-minion:
	@${MINION_APP} >/dev/null 2>&1 & echo $$! > $(PID_FILE)
	@cat $(PID_FILE) | sed "/^/s/^/PID: /"

stop-minion:
	@-touch $(PID_FILE)
	@-kill `cat $(PID_FILE)` 2> /dev/null || true
	@-rm $(PID_FILE)

start-shell:
	@${SHELL_APP}

start:
	@bash -c "$(MAKE) clean compile start-minion start-shell stop-minion"

docker-minion:
    @docker build . -f Minion.Dockerfile -t diligent-minion:latest

docker-boss:
    @docker build . -f Boss.Dockerfile -t diligent-boss:latest

docker-shell:
    @docker build . -f Shell.Dockerfile -t diligent-shell:latest

docker: docker-minion docker-shell docker-boss