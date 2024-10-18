"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[1415],{55888:(e,t,s)=>{s.r(t),s.d(t,{assets:()=>o,contentTitle:()=>i,default:()=>h,frontMatter:()=>a,metadata:()=>c,toc:()=>p});var n=s(74848),r=s(28453);const a={id:"user-purchasing-products",title:"OpenTelemetry Store - User Purchasing Products",description:"The OpenTelemetry Demo is an example application published by the OpenTelemtry CNCF project. This use case covers viewing recommended products before adding them to the shopping cart.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},i=void 0,c={id:"live-examples/opentelemetry-store/use-cases/user-purchasing-products",title:"OpenTelemetry Store - User Purchasing Products",description:"The OpenTelemetry Demo is an example application published by the OpenTelemtry CNCF project. This use case covers viewing recommended products before adding them to the shopping cart.",source:"@site/docs/live-examples/opentelemetry-store/use-cases/user-purchasing-products.mdx",sourceDirName:"live-examples/opentelemetry-store/use-cases",slug:"/live-examples/opentelemetry-store/use-cases/user-purchasing-products",permalink:"/live-examples/opentelemetry-store/use-cases/user-purchasing-products",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/live-examples/opentelemetry-store/use-cases/user-purchasing-products.mdx",tags:[],version:"current",frontMatter:{id:"user-purchasing-products",title:"OpenTelemetry Store - User Purchasing Products",description:"The OpenTelemetry Demo is an example application published by the OpenTelemtry CNCF project. This use case covers viewing recommended products before adding them to the shopping cart.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"guidesSidebar",previous:{title:"OpenTelemetry Store - Get recommended products",permalink:"/live-examples/opentelemetry-store/use-cases/get-recommended-products"},next:{title:"Pokeshop API Demo",permalink:"/live-examples/pokeshop/overview"}},o={},p=[{value:"Building a Test Suite for This Scenario",id:"building-a-test-suite-for-this-scenario",level:2},{value:"Mapping Environment Variables",id:"mapping-environment-variables",level:3},{value:"Creating Tests",id:"creating-tests",level:3},{value:"Creating the Test Suite",id:"creating-the-test-suite",level:3}];function d(e){const t={a:"a",code:"code",h2:"h2",h3:"h3",li:"li",ol:"ol",p:"p",pre:"pre",...(0,r.R)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(t.p,{children:"In this use case, we want to validate the following story:"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{children:"As a consumer, after landing at home page\nI want to see the shop recommended products, add the first one to my cart and pay for it on checkout\nSo I can have it shipped to my home\n"})}),"\n",(0,n.jsx)(t.p,{children:"Something interesting about this process is that it is a composition of many of the previous use cases, executed in sequence:"}),"\n",(0,n.jsxs)(t.ol,{children:["\n",(0,n.jsx)(t.li,{children:(0,n.jsx)(t.a,{href:"/live-examples/opentelemetry-store/use-cases/get-recommended-products",children:"Get Recommended Products"})}),"\n",(0,n.jsx)(t.li,{children:(0,n.jsx)(t.a,{href:"/live-examples/opentelemetry-store/use-cases/add-item-into-shopping-cart",children:"Add Item into Shopping Cart"})}),"\n",(0,n.jsx)(t.li,{children:(0,n.jsx)(t.a,{href:"/live-examples/opentelemetry-store/use-cases/check-shopping-cart-contents",children:"Check Shopping Cart Contents"})}),"\n",(0,n.jsx)(t.li,{children:(0,n.jsx)(t.a,{href:"/live-examples/opentelemetry-store/use-cases/checkout",children:"Checkout"})}),"\n"]}),"\n",(0,n.jsx)(t.p,{children:"So in this case, we need to trigger four tests in sequence to achieve test the entire scenario and make these tests share data."}),"\n",(0,n.jsx)(t.h2,{id:"building-a-test-suite-for-this-scenario",children:"Building a Test Suite for This Scenario"}),"\n",(0,n.jsxs)(t.p,{children:["Using Tracetest, we can do that by ",(0,n.jsx)(t.a,{href:"/web-ui/creating-tests",children:"creating a test"})," for each step and later grouping these tests as ",(0,n.jsx)(t.a,{href:"/web-ui/creating-test-suites",children:"Test Suites"})," that have an ",(0,n.jsx)(t.a,{href:"/concepts/variable-sets",children:"variable set"}),"."]}),"\n",(0,n.jsxs)(t.p,{children:["We can do that by creating the tests and Test Suites through the Web UI or using the CLI. In this example, we will use the CLI to create a Variable Set and then create the Test Suite with all tests needed. The ",(0,n.jsx)(t.a,{href:"/concepts/assertions",children:"assertions"})," that we will check are the same for every single test."]}),"\n",(0,n.jsx)(t.h3,{id:"mapping-environment-variables",children:"Mapping Environment Variables"}),"\n",(0,n.jsxs)(t.p,{children:["The first thing that we need to think about is to map the variables that are needed in this process. At first glance, we can identify the vars to the API address and the user ID:\nWith these variables, we can create the following definition file as saving as ",(0,n.jsx)(t.code,{children:"user-buying-products.env"}),":"]}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-sh",children:"OTEL_API_URL=http://otel-shop-demo-frontend:8080/api\nUSER_ID=2491f868-88f1-4345-8836-d5d8511a9f83\n"})}),"\n",(0,n.jsx)(t.h3,{id:"creating-tests",children:"Creating Tests"}),"\n",(0,n.jsxs)(t.p,{children:["After creating the environment file, we will create a test for each step, starting with ",(0,n.jsx)(t.a,{href:"/live-examples/opentelemetry-store/use-cases/get-recommended-products",children:"Get Recommended Products"}),", which will be saved as ",(0,n.jsx)(t.code,{children:"get-recommended-products.yaml"}),":"]}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-yaml",children:'type: Test\nspec:\n  name: Get recommended products\n  trigger:\n    type: http\n    httpRequest:\n      url: ${var:OTEL_API_URL}/recommendations?productIds=&sessionId=${var:USER_ID}&currencyCode=\n      method: GET\n      headers:\n      - key: Content-Type\n        value: application/json\n  specs:\n  - selector: span[tracetest.span.type="rpc" name="grpc.hipstershop.ProductCatalogService/GetProduct" rpc.system="grpc" rpc.method="GetProduct" rpc.service="hipstershop.ProductCatalogService"]\n    assertions: # It should have 4 products on this list.\n    - attr:tracetest.selected_spans.count = 4\n  - selector: span[tracetest.span.type="rpc" name="/hipstershop.FeatureFlagService/GetFlag" rpc.system="grpc" rpc.method="GetFlag" rpc.service="hipstershop.FeatureFlagService"]\n    assertions: # The feature flagger should be called for one product.\n    - attr:tracetest.selected_spans.count = 1\n  outputs:\n  - name: PRODUCT_ID\n    selector: span[tracetest.span.type="general" name="Tracetest trigger"]\n    value: attr:tracetest.response.body | json_path \'$[0].id\'\n'})}),"\n",(0,n.jsxs)(t.p,{children:["Note that we have one important changes here: we are now using environment variables on the definition, like ",(0,n.jsx)(t.code,{children:"${var:OTEL_API_URL}"})," and ",(0,n.jsx)(t.code,{children:"${var:USER_ID}"})," on the trigger section and an output to fetch the first ",(0,n.jsx)(t.code,{children:"${var:PRODUCT_ID}"})," that the user chose. This new environment variable will be used in the next tests."]}),"\n",(0,n.jsxs)(t.p,{children:["The next step is to define the ",(0,n.jsx)(t.a,{href:"/live-examples/opentelemetry-store/use-cases/add-item-into-shopping-cart",children:"Add Item into Shopping Cart"})," test, which will be saved as ",(0,n.jsx)(t.code,{children:"add-product-into-shopping-cart.yaml"}),":"]}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-yaml",children:'type: Test\nspec:\n  name: Add product into shopping cart\n  description: Add a selected product to user shopping cart\n  trigger:\n    type: http\n    httpRequest:\n      url: ${var:OTEL_API_URL}/cart\n      method: POST\n      headers:\n      - key: Content-Type\n        value: application/json\n      body: \'{"item":{"productId":"${var:PRODUCT_ID}","quantity":1},"userId":"${var:USER_ID}"}\'\n  specs:\n  - selector: span[tracetest.span.type="http" name="hipstershop.CartService/AddItem"]\n    # The correct ProductID was sent to the Product Catalog API.\n    assertions:\n    - attr:app.product.id = "${var:PRODUCT_ID}"\n  - selector: span[tracetest.span.type="database" name="HMSET" db.system="redis" db.redis.database_index="0"]\n    # The product persisted correctly on the shopping cart.\n    assertions:\n    - attr:tracetest.selected_spans.count >= 1\n'})}),"\n",(0,n.jsxs)(t.p,{children:["After that, we will ",(0,n.jsx)(t.a,{href:"/live-examples/opentelemetry-store/use-cases/check-shopping-cart-contents",children:"Check Shopping Cart Contents"})," (on ",(0,n.jsx)(t.code,{children:"check-shopping-cart-contents.yaml"}),"), simulating a user validating the products selected before finishing the purchase:"]}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-yaml",children:'type: Test\nspec:\n  name: Check shopping cart contents\n  trigger:\n    type: http\n    httpRequest:\n      url: ${var:OTEL_API_URL}/cart?sessionId=${var:USER_ID}&currencyCode=\n      method: GET\n      headers:\n      - key: Content-Type\n        value: application/json\n  specs:\n  - selector: span[tracetest.span.type="rpc" name="hipstershop.ProductCatalogService/GetProduct" rpc.system="grpc" rpc.method="GetProduct" rpc.service="hipstershop.ProductCatalogService"]\n    # The product previously added exists in the cart.\n    assertions:\n    - attr:app.product.id = "${var:PRODUCT_ID}"\n  - selector: span[tracetest.span.type="general" name="Tracetest trigger"]\n    # The size of the shopping cart should be at least 1.\n    assertions:\n    - attr:tracetest.response.body | json_path \'$.items.length\' >= 1\n'})}),"\n",(0,n.jsxs)(t.p,{children:["And finally, we have the ",(0,n.jsx)(t.a,{href:"/live-examples/opentelemetry-store/use-cases/checkout",children:"Checkout"})," action (",(0,n.jsx)(t.code,{children:"checkout.yaml"}),"), where the user inputs the billing and shipping info and finishes buying the item in the shopping cart:"]}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-yaml",children:'type: Test\nspec:\n  name: Checking out shopping cart\n  description: Checking out shopping cart\n  trigger:\n    type: http\n    httpRequest:\n      url: ${var:OTEL_API_URL}/checkout\n      method: POST\n      headers:\n      - key: Content-Type\n        value: application/json\n      body: \'{"userId":"${var:USER_ID}","email":"someone@example.com","address":{"streetAddress":"1600 Amphitheatre Parkway","state":"CA","country":"United States","city":"Mountain View","zipCode":"94043"},"userCurrency":"USD","creditCard":{"creditCardCvv":672,"creditCardExpirationMonth":1,"creditCardExpirationYear":2030,"creditCardNumber":"4432-8015-6152-0454"}}\'\n  specs:\n  - selector: span[tracetest.span.type="rpc" name="hipstershop.CheckoutService/PlaceOrder"\n      rpc.system="grpc" rpc.method="PlaceOrder" rpc.service="hipstershop.CheckoutService"]\n    assertions: \n    # An order was placed.\n    - attr:app.user.id = "${var:USER_ID}"\n    - attr:app.order.items.count = 1\n  - selector: span[tracetest.span.type="rpc" name="hipstershop.PaymentService/Charge" rpc.system="grpc" rpc.method="Charge" rpc.service="hipstershop.PaymentService"]\n    assertions: \n    # The user was charged.\n    - attr:rpc.grpc.status_code  =  0\n    - attr:tracetest.selected_spans.count >= 1\n  - selector: span[tracetest.span.type="rpc" name="hipstershop.ShippingService/ShipOrder" rpc.system="grpc" rpc.method="ShipOrder" rpc.service="hipstershop.ShippingService"]\n    assertions: \n    # The product was shipped.\n    - attr:rpc.grpc.status_code = 0\n    - attr:tracetest.selected_spans.count >= 1\n  - selector: span[tracetest.span.type="rpc" name="hipstershop.CartService/EmptyCart"\n      rpc.system="grpc" rpc.method="EmptyCart" rpc.service="hipstershop.CartService"]\n    assertions: \n    # The shopping cart was emptied.\n    - attr:rpc.grpc.status_code = 0\n    - attr:tracetest.selected_spans.count >= 1\n'})}),"\n",(0,n.jsx)(t.h3,{id:"creating-the-test-suite",children:"Creating the Test Suite"}),"\n",(0,n.jsxs)(t.p,{children:["Now we wrap these files and create a Test Suite that will run these tests in sequence and will fail if any of the tests fail. We will call it ",(0,n.jsx)(t.code,{children:"testsuite.yaml"}),":"]}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-yml",children:"type: TestSuite\nspec:\n  name: User purchasing products\n  description: Simulate a process of a user purchasing products on Astronomy store\n  steps:\n  - ./get-recommended-products.yaml\n  - ./add-product-into-shopping-cart.yaml\n  - ./check-shopping-cart-contents.yaml\n  - ./checkout.yaml\n"})}),"\n",(0,n.jsx)(t.p,{children:"By having the test, Test Suite and environment files in the same directory, we can call the CLI and execute this Test Suite:"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-sh",children:"tracetest run testsuite -f testsuite.yaml -e user-buying-products.env\n"})}),"\n",(0,n.jsx)(t.p,{children:"The result should be an output like this:"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-sh",children:"\u2714 User purchasing products (http://localhost:11633/testsuite/kRDUir0VR/run/1)\n        \u2714 Get recommended products (http://localhost:11633/test/XxH8irA4R/run/1/test)\n        \u2714 Add product into shopping cart (http://localhost:11633/test/j_N8i9AVR/run/1/test)\n        \u2714 Check shopping cart contents (http://localhost:11633/test/Y2jim9AVg/run/1/test)\n        \u2714 Checking out shopping cart (http://localhost:11633/test/VPCim90Vg/run/1/test)\n"})})]})}function h(e={}){const{wrapper:t}={...(0,r.R)(),...e.components};return t?(0,n.jsx)(t,{...e,children:(0,n.jsx)(d,{...e})}):d(e)}},28453:(e,t,s)=>{s.d(t,{R:()=>i,x:()=>c});var n=s(96540);const r={},a=n.createContext(r);function i(e){const t=n.useContext(a);return n.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function c(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:i(e.components),n.createElement(a.Provider,{value:t},e.children)}}}]);