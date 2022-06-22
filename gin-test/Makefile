SHELL=/bin/bash
.PHONY: build run clean

PROJECT_NAME:=$(shell echo $$(pwd) | awk -F '/' '{print $$NF}')

APPNAME:=$(if $(APPNAME),$(APPNAME),$(PROJECT_NAME))


build:
	@echo $(PROJECT_NAME)
	@go mod tidy
	@go mod vendor
	@go build -mod vendor -o $(APPNAME) $(shell pwd)/main.go
	@chmod +x $(APPNAME)
	@echo "构建成功$(APPNAME)"

run:
	./wait.bash | tee -a ./bash.log
	[ $? -gt 0 ] && { rm -rf $(APPNAME); killall $(APPNAME) ;}


clean:
	@rm -rf $(APPNAME)
	@killall $(APPNAME)