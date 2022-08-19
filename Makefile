# Configure targets
.PHONY: build-all build-boss build-minion build-shell clean test docker
default: build-all

BUILD_DIR:=./build
BOSS_APP:=$(BUILD_DIR)/boss
MINION_APP:=$(BUILD_DIR)/minion
SHELL_APP:=$(BUILD_DIR)/shell

APP_VERSION:=0.0.0
COMMIT_HASH:=$(shell git rev-parse --short HEAD 2>/dev/null || echo 'unknown')
GO_VERSION:=$(shell go version)
BUILD_TIME:=$(shell date)
LDFLAGS_COMMON:=$(LDFLAGS_COMMON) -X \"github.com/flipkart-incubator/diligent/pkg/buildinfo.AppVersion=$(APP_VERSION)\"
LDFLAGS_COMMON:=$(LDFLAGS_COMMON) -X \"github.com/flipkart-incubator/diligent/pkg/buildinfo.CommitHash=$(COMMIT_HASH)\"
LDFLAGS_COMMON:=$(LDFLAGS_COMMON) -X \"github.com/flipkart-incubator/diligent/pkg/buildinfo.GoVersion=$(GO_VERSION)\"
LDFLAGS_COMMON:=$(LDFLAGS_COMMON) -X \"github.com/flipkart-incubator/diligent/pkg/buildinfo.BuildTime=$(BUILD_TIME)\"

build-all: build-boss build-minion build-shell

build-boss:
	$(eval LDFLAGS_AV:=-X \"github.com/flipkart-incubator/diligent/pkg/buildinfo.AppName=diligent-boss\")
	$(eval LDFLAGS:=$(LDFLAGS_COMMON) $(LDFLAGS_AV))
	@echo "Building boss..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BOSS_APP) -ldflags="$(LDFLAGS)" github.com/flipkart-incubator/diligent/apps/boss
	@echo "Done"

build-minion:
	$(eval LDFLAGS_AV:=-X \"github.com/flipkart-incubator/diligent/pkg/buildinfo.AppName=diligent-minion\")
	$(eval LDFLAGS:=$(LDFLAGS_COMMON) $(LDFLAGS_AV))
	@echo "Building minion..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(MINION_APP) -ldflags="$(LDFLAGS)" github.com/flipkart-incubator/diligent/apps/minion
	@echo "Done"

build-shell:
	$(eval LDFLAGS_AV:=-X \"github.com/flipkart-incubator/diligent/pkg/buildinfo.AppName=diligent-shell\")
	$(eval LDFLAGS:=$(LDFLAGS_COMMON) $(LDFLAGS_AV))
	@echo "Building shell..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(SHELL_APP) -ldflags="$(LDFLAGSN)" github.com/flipkart-incubator/diligent/apps/shell
	@echo "Done"

clean:
	@echo "Cleaning..."
	@-rm -f $(BOSS_APP) $(MINION_APP) $(SHELL_APP)
	@-rmdir $(BUILD_DIR) 2> /dev/null
	@echo "Done"

test:
	@echo "Testing..."
	@go test ./...
	@echo "Done"

# docker-minion:
#     @docker build . -f Minion.Dockerfile -t diligent-minion:latest
#
# docker-boss:
#     @docker build . -f Boss.Dockerfile -t diligent-boss:latest
#
# docker-shell:
#     @docker build . -f Shell.Dockerfile -t diligent-shell:latest
#
# docker: docker-minion docker-shell docker-boss