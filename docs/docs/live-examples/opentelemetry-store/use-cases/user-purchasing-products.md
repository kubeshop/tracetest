# OpenTelemetry Store - User Purchasing Products

In this use case, we want to validate the following story:

```
As a consumer, after landing at home page
I want to see the shop recommended products, add the first one to my cart and pay for it on checkout
So I can have it shipped to my home
```

Something interesting about this process is that it is a composition of many of the previous use cases, executed in sequence:
1. [Get Recommended Products](./get-recommended-products.md)
2. [Add Item into Shopping Cart](./add-item-into-shopping-cart.md)
3. [Check Shopping Cart Contents](./check-shopping-cart-contents.md)
4. [Checkout](./checkout.md)

So in this case, we need to trigger four tests in sequence to achieve test the entire scenario and make these tests share data.

## Building a Transaction for This Scenario

Using Tracetest, we can do that by [creating a test](../../../web-ui/creating-tests.md) for each step and later grouping these tests as [transactions](../../../web-ui/creating-transactions.md) that have an [variable set](../../../concepts/variable-sets.md)](../../../concepts/variable-sets.md).
 
We can do that by creating the tests and transactions through the Web UI or using the CLI. In this example, we will use the CLI to create a Variable Set and then create the transaction with all tests needed. The [assertions](../../../concepts/assertions.md) that we will check are the same for every single test.

### Mapping Environment Variables 

The first thing that we need to think about is to map the variables that are needed in this process. At first glance, we can identify the vars to the API address and the user ID:
With these variables, we can create the following definition file as saving as `user-buying-products.env`:

```sh
OTEL_API_URL=http://otel-shop-demo-frontend:8080/api
USER_ID=2491f868-88f1-4345-8836-d5d8511a9f83
```

### Creating Tests

After creating the environment file, we will create a test for each step, starting with [Get Recommended Products](./get-recommended-products.md), which will be saved as `get-recommended-products.yaml`:

```yaml
type: Test
spec:
  name: Get recommended products
  trigger:
    type: http
    httpRequest:
      url: ${env:OTEL_API_URL}/recommendations?productIds=&sessionId=${env:USER_ID}&currencyCode=
      method: GET
      headers:
      - key: Content-Type
        value: application/json
  specs:
  - selector: span[tracetest.span.type="rpc" name="grpc.hipstershop.ProductCatalogService/GetProduct" rpc.system="grpc" rpc.method="GetProduct" rpc.service="hipstershop.ProductCatalogService"]
    assertions: # It should have 4 products on this list.
    - attr:tracetest.selected_spans.count = 4
  - selector: span[tracetest.span.type="rpc" name="/hipstershop.FeatureFlagService/GetFlag" rpc.system="grpc" rpc.method="GetFlag" rpc.service="hipstershop.FeatureFlagService"]
    assertions: # The feature flagger should be called for one product.
    - attr:tracetest.selected_spans.count = 1
  outputs:
  - name: PRODUCT_ID
    selector: span[tracetest.span.type="general" name="Tracetest trigger"]
    value: attr:tracetest.response.body | json_path '$[0].id'
```

Note that we have one important changes here: we are now using environment variables on the definition, like `${env:OTEL_API_URL}` and `${env:USER_ID}` on the trigger section and an output to fetch the first `${env:PRODUCT_ID}` that the user chose. This new environment variable will be used in the next tests.

The next step is to define the [Add Item into Shopping Cart](./add-item-into-shopping-cart.md) test, which will be saved as `add-product-into-shopping-cart.yaml`:

```yaml
type: Test
spec:
  name: Add product into shopping cart
  description: Add a selected product to user shopping cart
  trigger:
    type: http
    httpRequest:
      url: ${env:OTEL_API_URL}/cart
      method: POST
      headers:
      - key: Content-Type
        value: application/json
      body: '{"item":{"productId":"${env:PRODUCT_ID}","quantity":1},"userId":"${env:USER_ID}"}'
  specs:
  - selector: span[tracetest.span.type="http" name="hipstershop.CartService/AddItem"]
    # The correct ProductID was sent to the Product Catalog API.
    assertions:
    - attr:app.product.id = "${env:PRODUCT_ID}"
  - selector: span[tracetest.span.type="database" name="HMSET" db.system="redis" db.redis.database_index="0"]
    # The product persisted correctly on the shopping cart.
    assertions:
    - attr:tracetest.selected_spans.count >= 1
```

After that, we will [Check Shopping Cart Contents](./check-shopping-cart-contents.md) (on `check-shopping-cart-contents.yaml`), simulating a user validating the products selected before finishing the purchase:

```yaml
type: Test
spec:
  name: Check shopping cart contents
  trigger:
    type: http
    httpRequest:
      url: ${env:OTEL_API_URL}/cart?sessionId=${env:USER_ID}&currencyCode=
      method: GET
      headers:
      - key: Content-Type
        value: application/json
  specs:
  - selector: span[tracetest.span.type="rpc" name="hipstershop.ProductCatalogService/GetProduct" rpc.system="grpc" rpc.method="GetProduct" rpc.service="hipstershop.ProductCatalogService"]
    # The product previously added exists in the cart.
    assertions:
    - attr:app.product.id = "${env:PRODUCT_ID}"
  - selector: span[tracetest.span.type="general" name="Tracetest trigger"]
    # The size of the shopping cart should be at least 1.
    assertions:
    - attr:tracetest.response.body | json_path '$.items.length' >= 1
```

And finally, we have the [Checkout](./checkout.md) action (`checkout.yaml`), where the user inputs the billing and shipping info and finishes buying the item in the shopping cart:

```yaml
type: Test
spec:
  name: Checking out shopping cart
  description: Checking out shopping cart
  trigger:
    type: http
    httpRequest:
      url: ${env:OTEL_API_URL}/checkout
      method: POST
      headers:
      - key: Content-Type
        value: application/json
      body: '{"userId":"${env:USER_ID}","email":"someone@example.com","address":{"streetAddress":"1600 Amphitheatre Parkway","state":"CA","country":"United States","city":"Mountain View","zipCode":"94043"},"userCurrency":"USD","creditCard":{"creditCardCvv":672,"creditCardExpirationMonth":1,"creditCardExpirationYear":2030,"creditCardNumber":"4432-8015-6152-0454"}}'
  specs:
  - selector: span[tracetest.span.type="rpc" name="hipstershop.CheckoutService/PlaceOrder"
      rpc.system="grpc" rpc.method="PlaceOrder" rpc.service="hipstershop.CheckoutService"]
    assertions: 
    # An order was placed.
    - attr:app.user.id = "${env:USER_ID}"
    - attr:app.order.items.count = 1
  - selector: span[tracetest.span.type="rpc" name="hipstershop.PaymentService/Charge" rpc.system="grpc" rpc.method="Charge" rpc.service="hipstershop.PaymentService"]
    assertions: 
    # The user was charged.
    - attr:rpc.grpc.status_code  =  0
    - attr:tracetest.selected_spans.count >= 1
  - selector: span[tracetest.span.type="rpc" name="hipstershop.ShippingService/ShipOrder" rpc.system="grpc" rpc.method="ShipOrder" rpc.service="hipstershop.ShippingService"]
    assertions: 
    # The product was shipped.
    - attr:rpc.grpc.status_code = 0
    - attr:tracetest.selected_spans.count >= 1
  - selector: span[tracetest.span.type="rpc" name="hipstershop.CartService/EmptyCart"
      rpc.system="grpc" rpc.method="EmptyCart" rpc.service="hipstershop.CartService"]
    assertions: 
    # The shopping cart was emptied.
    - attr:rpc.grpc.status_code = 0
    - attr:tracetest.selected_spans.count >= 1
```

### Creating the Transaction

Now we wrap these files and create a transaction that will run these tests in sequence and will fail if any of the tests fail. We will call it `transaction.yaml`:

```yml
type: Transaction
spec:
  name: User purchasing products
  description: Simulate a process of a user purchasing products on Astronomy store
  steps:
  - ./get-recommended-products.yaml
  - ./add-product-into-shopping-cart.yaml
  - ./check-shopping-cart-contents.yaml
  - ./checkout.yaml
```

By having the test, transaction and environment files in the same directory, we can call the CLI and execute this transaction:

```sh
tracetest run transaction -f transaction.yaml -e user-buying-products.env
```

The result should be an output like this:

```sh
✔ User purchasing products (http://localhost:11633/transaction/kRDUir0VR/run/1)
        ✔ Get recommended products (http://localhost:11633/test/XxH8irA4R/run/1/test)
        ✔ Add product into shopping cart (http://localhost:11633/test/j_N8i9AVR/run/1/test)
        ✔ Check shopping cart contents (http://localhost:11633/test/Y2jim9AVg/run/1/test)
        ✔ Checking out shopping cart (http://localhost:11633/test/VPCim90Vg/run/1/test)
```
