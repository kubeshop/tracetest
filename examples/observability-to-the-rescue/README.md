# Observability to the Rescue! support material

To run this example just run:
```sh
docker compose up
```

Example of request that runs correctly (valid payment, without risk analysis):
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

Example of request that runs correctly (valid payment, risk analysis):
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

Example of request that runs correctly (denied payment):
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

Example of request that has error:
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
