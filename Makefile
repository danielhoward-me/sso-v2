# You may need to install the following dependencies:
# - jet (go install github.com/go-jet/jet/v2/cmd/jet@latest)
# - concurrently (npm i -g concurrently)
# - wait-on (npm i -g wait-on)
# - nodemon (npm i -g nodemon)

# Backend commands
backend-jet:
	jet -dsn=postgresql://sso:sso@local.danielhoward.me:5432/sso?sslmode=disable -schema=public -path=./.gen
	echo "Moving generated files to internal/db/schema"
	rm -rf ./backend/internal/db/schema
	mv ./.gen/sso/public ./backend/internal/db/schema
	rm -rf ./.gen
backend-dev:
	export PORT=3001; \
	export GIN_MODE=debug; \
	export PGUSER=sso; \
	export PGPASSWORD=sso; \
	export PGHOST=local.danielhoward.me; \
	export PGPORT=5432; \
	export PGDATABASE=sso; \
	wait-on tcp:$$PGPORT && nodemon --ext "go" --watch "./backend/**/*.*" --exec "go run ./backend"

# Frontend commands


# Root commands
dev-docker:
	docker compose -f ./dev/compose.dev.yml -p sso-dev up
dev:
	concurrently -k -c yellow,blue -n BACKEND,DOCKER "make backend-dev" "make dev-db-proxy"

.PHONY: backend-jet backend-dev dev-docker dev
