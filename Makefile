run:
	docker compose up
build:
	docker compose build
stop:
	docker compose down
restart:
	docker compose restart
demo-otel-run:
	docker compose -f docker-compose.yaml -f ./local-config/docker-compose.otel-demo.yaml --env-file ./local-config/otel-demo.env up
demo-otel-stop:
	docker compose -f docker-compose.yaml -f ./local-config/docker-compose.otel-demo.yaml --env-file ./local-config/otel-demo.env down
