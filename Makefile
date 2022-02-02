
server-generate:
	openapi-generator-cli generate -i api/openapi.yaml -g go-server -o server/

	