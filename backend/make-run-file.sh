cd "$(dirname "$0")"

CMD=$(cat <<-END
echo "WARNING:" &&
echo " - Starting image with sleep command due to missing generator files" &&
echo " - Please keep this running and run \\\`make backend-generate\\\` in another terminal" &&
echo " - Then restart this one with \\\`make dev-build\\\`" &&
sleep 999
END
)

if [ -d "./internal/db/schema" ] &&
    [ -f "./internal/server/types.gen.go" ] &&
    [ -f "./internal/server/accountapi/api.gen.go" ] &&
    [ -f "./internal/server/internalapi/api.gen.go" ] &&
    [ -f "./internal/server/ssoapi/api.gen.go" ]; then
    CMD="go run ./cmd/sso-backend"
fi

echo $CMD > ./run.sh
