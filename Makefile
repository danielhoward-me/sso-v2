dev: # Run the docker development environment
	docker compose -f ./dev/compose.dev.yml up --watch
dev-build: # Reuild the development containers and run the docker development environment
	docker compose -f ./dev/compose.dev.yml up --watch --build

backend-sh: # Run a shell in the backend container
	docker compose -f ./dev/compose.dev.yml exec backend sh
backend-generate: # Generate the backend auto generated code
	docker compose -f ./dev/compose.dev.yml exec backend go generate ./...

.PHONY: dev dev-build backend-sh backend-generate
