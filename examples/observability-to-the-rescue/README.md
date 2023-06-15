# Observability to the Rescue! support material

This is a support material for the presentation "Observability to the Rescue! Monitoring and testing APIs with OpenTelemetry".

To run this example just run:
```sh
docker compose up
```

### Requests that you can do this example

Valid payment without risk analysis scenario:
```sh
curl --location 'http://localhost:10013/executePaymentOrder' \
--header 'Content-Type: application/json' \
--data '{
    "walletId": 2,
    "yearsAsACustomer": 1
}'

# Output
# {
#     "status": "executed"
# }
```

Valid payment with risk analysis scenario:
```sh
curl --location 'http://localhost:10013/executePaymentOrder' \
--header 'Content-Type: application/json' \
--data '{
    "walletId": 4,
    "yearsAsACustomer": 1
}'

# Output
# {
#     "status": "executed"
# }
```

Denied payment scenario:
```sh
curl --location 'http://localhost:10013/executePaymentOrder' \
--header 'Content-Type: application/json' \
--data '{
    "walletId": 5,
    "yearsAsACustomer": 1
}'

# Output
# {
#     "status": "denied"
# }
```

Request with error scenario
```sh
curl --location 'http://localhost:10013/executePaymentOrder' \
--header 'Content-Type: application/json' \
--data '{
    "walletId": 4,
    "yearsAsACustomer": 0
}'

# Output
# internal error!
```

## Trace-based tests that you can run

```sh
tracetest test run -w -d ./tracetest/tests/test-with-error.yaml
```

```sh
tracetest test run -w -d ./tracetest/tests/test-with-success.yaml
```
