# Go module env variables
export GOFLAGS		= -mod=vendor
export GO111MODULE	= on

APP_NAME 		:= demo-server
TARGET			:= target/bin/${APP_NAME}
APP_URL			:= http://127.0.0.1:8000/random/map
.DEFAULT_GOAL	:= help

.PHONY: help
help:  ## Display help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"}     /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-16s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

.PHONY: run
run: ## Run app with 1 CPU
	env GOMAXPROCS=1 go run main.go

.PHONY: build
build: ## Build binary
	env GOOS=linux CGO_ENABLED=0 go build -o ${TARGET} main.go

.PHONY: test
test: ## Run tests
	go test -race ./...

.PHONY: benchmark
benchmark: ## Run benchmarks and create report
	mkdir -p target/bench
	go test -bench=. -benchmem -benchtime=5s ./... | tee target/bench/$(shell date +%Y-%m-%d_%H:%M:%S).txt

.PHONY: benchcmp
benchcmp: ## Run benchmark and compare to previous
	benchcmp target/bench/$(shell ls --sort time target/bench/ | sed -n 2p) target/bench/$(shell ls --sort time target/bench/ | sed -n 1p)

.PHONY: clean
clean: ## Clean target folder
	rm -r target/*

.PHONY: wrk
wrk: ## Run wrk against app
	mkdir -p target/wrk
	wrk -d15s ${APP_URL} | tee target/wrk/$(shell date +%Y-%m-%d_%H:%M:%S).txt

.PHONY: vegeta
vegeta: ## Run vegeta against app
	echo "GET ${APP_URL}" | vegeta attack -name go -duration 15s -rate 2000 | tee target/vegeta_result.bin | vegeta report; vegeta plot target/vegeta_result.bin > target/plot.html

.PHONY: cpu-profile
cpu-profile: ## Get cpu profile
	wget -O target/cpu.prof 'http://127.0.0.1:8000/debug/pprof/profile?seconds=10'

.PHONY: heap-profile
heap-profile: ## Get heap profile
	wget -O target/mem.prof 'http://127.0.0.1:8000/debug/pprof/heap'

trace-profile: ## Get trace profile
	wget -O target/trace.prof 'http://127.0.0.1:8000/debug/pprof/trace?seconds=2'

.PHONY: escape-info
escape-info: ## Get escape info from compiler
	go build -o /dev/null -gcflags='-m' main.go

.PHONY: bound-info
bound-info: ## Get bound check info from compiler
	go build -o /dev/null -gcflags='-d=ssa/check_bce/debug=1' main.go

.PHONY: cache-info
cache-info: ## Get escape info from compiler
	perf c2c record -u -p $(shell fuser 8000/tcp 2>/dev/null) --call-graph dwarf sleep 10
	perf c2c report -g --stdio
