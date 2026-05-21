# operator-commons Makefile — RFC-0017 §3.3 표준 타겟 (라이브러리 특성).
# 본 라이브러리는 cmd/main.go 부재 → build/run/deploy 타겟 없음.
# 필수 타겟: lint, test, audit (validate, sbom 제외).

SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

.PHONY: all
all: lint test ## 기본 타겟 (CI 동등)

.PHONY: help
help: ## 본 Makefile 의 사용 가능 타겟 목록
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make \033[36m<target>\033[0m\n\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

##@ Quality

.PHONY: fmt
fmt: ## go fmt (모든 패키지)
	go fmt ./...

.PHONY: vet
vet: ## go vet (모든 패키지)
	go vet ./...

.PHONY: lint
lint: ## golangci-lint 풀스캔 (.golangci.yml + .custom-gcl.yml logcheck plugin)
	@command -v golangci-lint >/dev/null || { echo "[error] golangci-lint 미설치 — https://golangci-lint.run/usage/install/" >&2; exit 1; }
	golangci-lint run --config .golangci.yml ./...

.PHONY: test
test: ## go test (race + coverage profile cover.out)
	go test -race -count=1 -timeout=120s -coverprofile=cover.out ./...

.PHONY: cover
cover: test ## test 결과를 HTML 로 변환 (cover.html)
	go tool cover -html=cover.out -o cover.html
	@echo "→ cover.html"

##@ Security / Dependency

.PHONY: audit
audit: ## govulncheck (Go module CVE call-graph 검사)
	@command -v govulncheck >/dev/null || { echo "[warn] govulncheck 미설치 — go install golang.org/x/vuln/cmd/govulncheck@latest" >&2; exit 0; }
	govulncheck ./...

.PHONY: audit-quality
audit-quality: ## 5 repo production-grade 자동 측정 (P0/P1/P2/OP/C 50+ 항목, ADR-0013)
	@bash scripts/audit-production-grade.sh

.PHONY: tidy
tidy: ## go mod tidy + diff 검사 (drift 차단)
	@cp go.mod /tmp/.gomod.bak; cp go.sum /tmp/.gosum.bak
	@go mod tidy
	@if ! diff -q go.mod /tmp/.gomod.bak >/dev/null || ! diff -q go.sum /tmp/.gosum.bak >/dev/null; then \
		echo "[info] go mod tidy 가 변경을 발생시킴 — 본 변경을 commit 하세요"; \
	else \
		echo "[ok] go.mod / go.sum drift 없음"; \
	fi
	@rm -f /tmp/.gomod.bak /tmp/.gosum.bak

##@ Release

.PHONY: tag
tag: ## annotated tag 생성 안내 (수동 — v0.X.Y 인자 필요)
	@echo "사용법: git tag -a v0.X.Y -m 'release v0.X.Y' && git push origin v0.X.Y"
	@echo "후속: 3 consumer operator 의 go.mod require 버전 bump PR 작성"

.PHONY: release
release: ## 자동 release pipeline (scripts/release.sh, ADR-0014). 사용: make release VERSION=v0.8.0
	@[ -n "$(VERSION)" ] || { echo "Usage: make release VERSION=v0.8.0"; exit 1; }
	bash scripts/release.sh $(VERSION)
