default: dev

dev: # Run the docker development environment
	docker compose -f ./dev/compose.dev.yml up --watch
dev-build: # Reuild the development containers and run the docker development environment
	docker compose -f ./dev/compose.dev.yml up --watch --build

backend-sh: # Run a shell in the backend container
	docker compose -f ./dev/compose.dev.yml exec backend sh
backend-generate: # Generate the backend auto generated code
	docker compose -f ./dev/compose.dev.yml exec backend go generate ./...
	docker compose -f ./dev/compose.dev.yml cp backend:/app/internal/api/accountapi/api.gen.go ./backend/internal/api/accountapi/api.gen.go
	docker compose -f ./dev/compose.dev.yml cp backend:/app/internal/api/ssoapi/api.gen.go ./backend/internal/api/ssoapi/api.gen.go
	docker compose -f ./dev/compose.dev.yml cp backend:/app/internal/api/ssointernalapi/api.gen.go ./backend/internal/api/ssointernalapi/api.gen.go
	rm -rf ./backend/internal/db/schema
	docker compose -f ./dev/compose.dev.yml cp backend:/app/internal/db/.gen/sso/public ./backend/internal/db
	mv ./backend/internal/db/public ./backend/internal/db/schema
	for file in $$(find ./backend/internal/db/schema -name '*.go'); do \
		mv "$$file" "$${file%.go}.gen.go"; \
	done

.PHONY: dev dev-build backend-sh backend-generate
