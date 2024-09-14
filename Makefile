backend-jet:
	cd backend && go run github.com/go-jet/jet/v2/cmd/jet -dsn=postgresql://sso:sso@local.danielhoward.me:5432/sso?sslmode=disable -schema=public -path=./.gen
	rm -rf ./backend/internal/db/schema
	mkdir -p ./backend/internal/db/schema
	mv ./backend/.gen/sso/public ./backend/internal/db/schema
	rm -rf ./backend/.gen
	for file in $$(find ./backend/internal/db/schema -name '*.go'); do \
		mv "$$file" "$${file%.go}.gen.go"; \
	done

dev:
	docker compose -f ./dev/compose.dev.yml up
dev-build:
	docker compose -f ./dev/compose.dev.yml up --build
generate:
	docker compose -f ./dev/compose.dev.yml exec backend sh

.PHONY: backend-jet dev dev-build
