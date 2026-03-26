.PHONY: dev generate check test build preview deploy reset

# ── Development ──────────────────────────────
dev:
	@trap 'kill 0' INT; \
	go run ./cmd/server & \
	(cd web && vp dev) & \
	wait

# ── Code generation ──────────────────────────
generate:
	sqlc generate -f db/sqlc.yaml
	oapi-codegen -generate models,chi-server,spec -package api -o internal/api/api.gen.go api/openapi.yaml
	cd web && npx openapi-typescript ../api/openapi.yaml -o src/lib/api/v1.d.ts && vp fmt src/lib/api/v1.d.ts

# ── Checks ───────────────────────────────────
check:
	go vet ./...
	cd web && vp check

# ── Tests ────────────────────────────────────
test:
	go test ./...

# ── Build ────────────────────────────────────
build:
	cd web && vp build
	go build -o ./bin/fabel ./cmd/server
	@echo "Built: $$(ls -lh ./bin/fabel | awk '{print $$5}')"

# ── Preview ──────────────────────────────────
preview: build
	./bin/fabel

# ── Deploy ───────────────────────────────────
deploy: build
	@echo "TODO: scp + systemctl restart"

# ── Utilities ────────────────────────────────
reset:
	rm -f fabel.db fabel.db-wal fabel.db-shm
