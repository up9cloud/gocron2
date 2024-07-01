GO111MODULE=on

.PHONY: build
build: gocron2 gocron2-node

.PHONY: build-race
build-race: enable-race build

.PHONY: gocron2
gocron2:
	go build $(RACE) -o bin/gocron2 ./cmd/gocron2/main.go

.PHONY: gocron2-node
gocron2-node:
	go build $(RACE) -o bin/gocron2-node ./cmd/gocron2-node/main.go

# 编译并运行
.PHONY: run
run: build kill
	./bin/gocron2-node -allow-root &
	./bin/gocron2 web -e dev

.PHONY: run-race
run-race: enable-race run

.PHONY: kill
kill:
	-killall gocron2-node

.PHONY: test
test:
	go test $(RACE) ./...

.PHONY: test-race
test-race: enable-race test

.PHONY: enable-race
enable-race:
	$(eval RACE = -race)

# 打包: 生成当前系统的压缩包檔案
.PHONY: package
package: install-vue build-vue statik
	bash ./package.sh

# 打包: 生成所有系統的压缩包
.PHONY: package-all
package-all: install-vue build-vue statik
	bash ./package.sh -p 'linux darwin windows' -a 'amd64 arm64'

.PHONY: install-vue
install-vue:
	cd web/vue && npm install

.PHONY: build-vue
build-vue:
	cd web/vue && npm run build
	cp -r web/vue/dist/* web/public/

.PHONY: run-vue
run-vue:
	cd web/vue && npm run dev

.PHONY: statik
statik:
	go install github.com/rakyll/statik
	go generate ./...

.PHONY: lint
	golangci-lint run

.PHONY: clean
clean:
	rm bin/gocron2
	rm bin/gocron2-node
