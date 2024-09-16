default: dev

dev: # Run the docker development environment
	docker compose -f ./dev/docker-compose.yaml up --watch
dev-build: # Reuild the development containers and run the docker development environment
	docker compose -f ./dev/docker-compose.yaml up --watch --build

backend-sh: # Run a shell in the backend container
	docker compose -f ./dev/docker-compose.yaml exec backend sh
backend-generate: # Generate the backend auto generated code
	docker compose -f ./dev/docker-compose.yaml exec backend go generate -x ./...
	docker compose -f ./dev/docker-compose.yaml cp backend:/app/internal/server/types.gen.go ./backend/internal/server/types.gen.go
	docker compose -f ./dev/docker-compose.yaml cp backend:/app/internal/server/accountapi/api.gen.go ./backend/internal/server/accountapi/api.gen.go
	docker compose -f ./dev/docker-compose.yaml cp backend:/app/internal/server/internalapi/api.gen.go ./backend/internal/server/internalapi/api.gen.go
	docker compose -f ./dev/docker-compose.yaml cp backend:/app/internal/server/ssoapi/api.gen.go ./backend/internal/server/ssoapi/api.gen.go
	rm -rf ./backend/internal/db/schema
	docker compose -f ./dev/docker-compose.yaml cp backend:/app/internal/db/.gen/sso/public ./backend/internal/db
	mv ./backend/internal/db/public ./backend/internal/db/schema
	for file in $$(find ./backend/internal/db/schema -name '*.go'); do \
		mv "$$file" "$${file%.go}.gen.go"; \
	done

.PHONY: dev dev-build backend-sh backend-generate
