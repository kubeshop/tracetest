# Advanced Selectors

If you find yourself in a position where you cannot select complex spans, you can use our advanced selectors to help in that task. It enables you selecting spans that are impossible to select using just basic selectors.

In order to present each selector feature as easily as possible, we will use a theoretical scenario of an e-commerce application.

The system we will inspect has this flow:

```mermaid
flowchart LR
    start((start))
    subgraph purchase
        cart-api
        purchase-api
    end
    subgraph auth
        auth-api
        auth-storage[(db)]
    end
    subgraph product
        product-api
        product-storage[(db)]
    end
    subgraph notification
        notification-api
        kafka{{kafka}}
        external-notification-service{{external service}}
    end

    start -->|1. Close order| cart-api
    cart-api-->|5. send buy action| purchase-api
    purchase-api --> |7. Send notification to user|notification-api
    purchase-api -->|6. can product be bought by user?| auth-api 
    auth-api --> auth-storage
    cart-api -->|2. is product available?| product-api
    product-api -->|4. can user view product?| auth-api
    product-api -->|3. retrieve product| product-storage

    notification-api -->|8| kafka
    kafka -->|9| external-notification-service
```

And it generates the following trace:


```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K
```

# Features

## Empty selector
By providing an empty selector, all spans from the trace are selected. Note that an empty selector is an empty string. Providing `span` or `span[]` as a selector will result as a syntax error.

## Filter by attributes
The most basic way of filtering the spans you want to apply an assertion on is by using its attributes. A good starting example would be filtering all spans of type `http`:

```css
span[tracetest.span.type="http"]
```

This would select the following spans:

```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    A:::selectedSpan
    B:::selectedSpan
    D:::selectedSpan
    F:::selectedSpan
    G:::selectedSpan
    I:::selectedSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef candidateSpan fill:#FF6905, color:#ffffff
```

### AND condition
If you need to narrow down your results, you can provide multiple properties in the selector by separating them usign a space. So, let's say we want all `http` spans **AND** were created by the service `cart-api`:

```css
span[tracetest.span.type="http" service.name="cart-api"]
```

This would select the following spans:

```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    A:::selectedSpan
    F:::selectedSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef candidateSpan fill:#FF6905, color:#ffffff
```

### OR condition
Sometimes we want to have a broader result by selecting spans that match different selectors. Let's say we have to get all spans from our services, but not from any other external service.

```css
span[service.name="api-product"], span[service.name="api-auth"], span[service.name="api-notification"], span[service.name="api-cart"]
```

This would select the following spans:

```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    A:::selectedSpan
    B:::selectedSpan
    C:::selectedSpan
    D:::selectedSpan
    E:::selectedSpan
    F:::selectedSpan
    G:::selectedSpan
    H:::selectedSpan
    I:::selectedSpan
    J:::selectedSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef candidateSpan fill:#FF6905, color:#ffffff
```

Each span selector will be executed individually and the results will be merged together, creating a list of all spans that match any of the provided span selectors.

### Contains operator
Although it is possible to filter several span selectors at once to get a broader result, it might become verbose very quickly. The previous example can be written in another way to reduce its complexity:

```css
span[service.name contains "api"]
```

This would select the same spans as the previous example:

```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    A:::selectedSpan
    B:::selectedSpan
    C:::selectedSpan
    D:::selectedSpan
    E:::selectedSpan
    F:::selectedSpan
    G:::selectedSpan
    H:::selectedSpan
    I:::selectedSpan
    J:::selectedSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef candidateSpan fill:#FF6905, color:#ffffff
```
## pseudo-classes support

Sometimes filtering by attributes is not enough because we might have two or three identical spans in the tree, but we only want to assert one of them. For example, imagine a system that has a `retry` policy for all the HTTP requests it sends. How would we allow a user to validate if the `third` execution was successful without asserting the other two spans?

This is where pseudo-classes enter the scene. Those are ways of filtering spans by data that is not present in the span itself. For example, the order which the span appears.

> :warning: Today we support only `first`, `last`, and `nth_child`. If you think we should implement another one, please open an issue and explain why it is important and how it should behave.

For the examples of the three pseudo-classes, let's consider we want to select a specific `http` span based on when it happens.

```css
span[tracetest.span.type="http"]
```

This will select the following spans:

```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    A:::selectedSpan
    B:::selectedSpan
    D:::selectedSpan
    F:::selectedSpan
    G:::selectedSpan
    I:::selectedSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef candidateSpan fill:#FF6905, color:#ffffff
```

### :first
This would return the first appearing span from the list

```css
span[tracetest.span.type="http"]:first
```

```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    A:::selectedSpan
    B:::candidateSpan
    D:::candidateSpan
    F:::candidateSpan
    G:::candidateSpan
    I:::candidateSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef candidateSpan fill:#FF6905, color:#ffffff
```

### :last
This would return the last appearing span from the list

```css
span[tracetest.span.type="http"]:last
```

```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    A:::candidateSpan
    B:::candidateSpan
    D:::candidateSpan
    F:::candidateSpan
    G:::candidateSpan
    I:::selectedSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef candidateSpan fill:#FF6905, color:#ffffff
```

### :nth_child
This enables you to fetch any item from the list based on its index. `n` starts at 1 (first element) and ends at `length` (last element). Any invalid `n` value will return in an empty list of spans being returned

```css
span[tracetest.span.type="http"]:nth_child(3)
```

```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    A:::candidateSpan
    B:::candidateSpan
    D:::selectedSpan
    F:::candidateSpan
    G:::candidateSpan
    I:::candidateSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef candidateSpan fill:#FF6905, color:#ffffff
```


## Parent-child relation filtering
Even with all those capabilities, we might have problems with ambiguous selectors returning several spans when just a few were intended.

In our example, you can notice that `auth-api` is called twice from different parts of the trace. At first by `product-api` and later by `cart-api`.

What if I want to test if a product only available in US can be bought in UK? The product can be seen by the user, but it cannot be bought if the user is outside US. Certainly I cannot apply the same assertions on all `auth-api` spans, otherwise the test will not pass.

> :information_source: When you filter by parent-child relation, it matches spans recursively in all levels bellow the parent. This doesn't match only direct children of the parent, but all other spans in the sub-tree.

For example:

```css
span[service.name="auth-api" tracetest.span.type="http"]
```

Will return:
```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    D:::selectedSpan
    G:::selectedSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef candidateSpan fill:#FF6905, color:#ffffff
```

This is a problem, because if we apply the same assertion to both spans, one of them will fail. We could try to use `nth_child` but that could break if a http request failed and the retry policy kicked in. Thus, the only way of filtering that is based on the context when it was generated. For example: using its parent span to do so.

We could use the `purchase products` parent to ensure just `http` class to the `auth-api` triggered by the `purchase-api` would be selected:

```css
span[service.name="cart-api", name="purchase products"] span[service.name="auth-api" tracetest.span.type="http"]
```

This would find the parent span and only select the spans that are descedents of that parent and match the provided filter:

```mermaid
flowchart TD
    A[" id: 1
        close order
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    B[" id: 2
        is product available

        attributes:
        service.name: product-api
        tracetest.span.type: http
        http.method: GET
    "]
    C[" id: 3
        get product information

        attributes:
        service.name: product-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    D[" id: 4
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    E[" id: 5
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    F[" id: 6
        purchase products
        
        attributes:
        service.name: cart-api
        tracetest.span.type: http
        http.method: POST
    "]
    G[" id: 7
        get user can access

        attributes:
        service.name: auth-api
        tracetest.span.type: http
        http.method: GET
    "]
    H[" id: 8
        get user auth information

        attributes:
        service.name: auth-api
        tracetest.span.type: db
        db.statement: SELECT * FROM ...
    "]
    I[" id: 9
        notify user

        attributes:
        service.name: notification-api
        tracetest.span.type: http
        http.method: POST
    "]
    J[" id: 10
        send message to kafka

        attributes:
        service.name: notification-api
        tracetest.span.type: messaging
    "]
    K[" id: 10
        send message to users

        attributes:
        service.name: external-service
        tracetest.span.type: messaging
    "]


    A -->|1| B
    B -->|2| C
    B -->|3| D
    D -->|4| E
    A -->|5| F
    F -->|6| G
    G -->|7| H
    F -->|8| I
    I -->|9| J
    J -->|10| K

    F:::parentSpan
    G:::selectedSpan

    classDef selectedSpan fill:#439846, color:#ffffff
    classDef parentSpan fill:#3792cb, color:#ffffff
```