type: Demo
spec:
  id: ""
  name: dev-updated
  enabled: true
  type: otelstore
  opentelemetryStore:
    frontendEndpoint: http://dev-updated-frontend:9000
    productCatalogEndpoint: http://dev-updated-product:8081
    cartEndpoint: http://dev-updated-cart:8082
    checkoutEndpoint: http://dev-updated-checkout:8083
