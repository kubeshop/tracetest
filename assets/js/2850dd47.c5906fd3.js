"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[6797],{95122:(e,s,n)=>{n.r(s),n.d(s,{assets:()=>c,contentTitle:()=>i,default:()=>d,frontMatter:()=>a,metadata:()=>r,toc:()=>p});var t=n(74848),o=n(28453);const a={id:"import-pokemon",title:"Pokeshop API - Import Pokemon",description:"As a testing ground, the Tracetest team has implemented a sample API instrumented with OpenTelemetry around the PokeAPI. This use case showcases a more complex scenario involving an async process.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},i=void 0,r={id:"live-examples/pokeshop/use-cases/import-pokemon",title:"Pokeshop API - Import Pokemon",description:"As a testing ground, the Tracetest team has implemented a sample API instrumented with OpenTelemetry around the PokeAPI. This use case showcases a more complex scenario involving an async process.",source:"@site/docs/live-examples/pokeshop/use-cases/import-pokemon.mdx",sourceDirName:"live-examples/pokeshop/use-cases",slug:"/live-examples/pokeshop/use-cases/import-pokemon",permalink:"/live-examples/pokeshop/use-cases/import-pokemon",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/live-examples/pokeshop/use-cases/import-pokemon.mdx",tags:[],version:"current",frontMatter:{id:"import-pokemon",title:"Pokeshop API - Import Pokemon",description:"As a testing ground, the Tracetest team has implemented a sample API instrumented with OpenTelemetry around the PokeAPI. This use case showcases a more complex scenario involving an async process.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"guidesSidebar",previous:{title:"Pokeshop API - Get Pokemon by ID",permalink:"/live-examples/pokeshop/use-cases/get-pokemon-by-id"},next:{title:"Pokeshop API - Import Pokemon from Stream",permalink:"/live-examples/pokeshop/use-cases/import-pokemon-from-stream"}},c={},p=[{value:"Building a Test for This Scenario",id:"building-a-test-for-this-scenario",level:2},{value:"Traces",id:"traces",level:3},{value:"Assertions",id:"assertions",level:3},{value:"Test Definition",id:"test-definition",level:3}];function l(e){const s={a:"a",code:"code",h2:"h2",h3:"h3",img:"img",li:"li",mermaid:"mermaid",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,o.R)(),...e.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(s.p,{children:"This use case showcases a more complex scenario involving an async process. Usually, when working with microservices, there are use cases where some of the processing needs to happen asynchronously, for example, when triggering a user notification, generating reports or processing a payment order. With this endpoint, we provide an example of how users can implement trace-based testing for such scenarios."}),"\n",(0,t.jsx)(s.p,{children:"Here the process is split into two phases:"}),"\n",(0,t.jsxs)(s.ol,{children:["\n",(0,t.jsx)(s.li,{children:"An API call that enqueues an import request to a queue."}),"\n"]}),"\n",(0,t.jsx)(s.mermaid,{value:'sequenceDiagram\n    participant Endpoint as POST /pokemon\n    participant API as API\n    participant Queue as RabbitMQ\n    \n    Endpoint->>API: request\n\n    alt request is invalid\n        API--\x3e>Endpoint: 400 Bad Request <br> <List of errors>\n    end\n\n    API->>Queue: enqueue "import" message\n    Queue--\x3e>API: message queued\n\n    API--\x3e>Endpoint: 201 Created'}),"\n",(0,t.jsxs)(s.ol,{start:"2",children:["\n",(0,t.jsx)(s.li,{children:"A Worker that dequeues messages and completes the async process."}),"\n"]}),"\n",(0,t.jsx)(s.mermaid,{value:'sequenceDiagram\n    participant Queue as RabbitMQ\n    participant Worker as Queue Worker\n    participant ExternalAPI as PokeAPI\n    participant Database as Postgres\n    \n    Queue->>Worker: dequeue "import" message\n\n    Worker->>ExternalAPI: get pokemon info\n    ExternalAPI--\x3e>Worker: pokemon info\n\n    Worker->>Database: save pokemon\n    Database--\x3e>Worker: pokemon saved\n  \n    alt if successful\n      Worker--\x3e>Queue: ack <br> queue can delete message\n    else is failed\n      Worker--\x3e>Queue: nack <br> queue keep message for retries\n    end'}),"\n",(0,t.jsxs)(s.p,{children:["You can trigger this use case by calling the endpoint ",(0,t.jsx)(s.code,{children:"POST /pokemon/import"}),", with the following request body:"]}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-json",children:'{\n  "id":  52\n}\n'})}),"\n",(0,t.jsx)(s.p,{children:"It should return the following payload:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-json",children:'{\n  "id":  52\n}\n'})}),"\n",(0,t.jsx)(s.h2,{id:"building-a-test-for-this-scenario",children:"Building a Test for This Scenario"}),"\n",(0,t.jsxs)(s.p,{children:["Using Tracetest, we can ",(0,t.jsx)(s.a,{href:"/web-ui/creating-tests",children:"create a test"})," that will execute an API call on ",(0,t.jsx)(s.code,{children:"POST /pokemon/import"})," and validate the following properties:"]}),"\n",(0,t.jsxs)(s.ul,{children:["\n",(0,t.jsx)(s.li,{children:"The API should enqueue an import task and return HTTP 200 OK."}),"\n",(0,t.jsx)(s.li,{children:"The worker should dequeue the import task."}),"\n",(0,t.jsx)(s.li,{children:"PokeAPI should return a valid response."}),"\n",(0,t.jsx)(s.li,{children:"The database should respond with low latency (< 200ms)."}),"\n"]}),"\n",(0,t.jsx)(s.h3,{id:"traces",children:"Traces"}),"\n",(0,t.jsxs)(s.p,{children:["Running these tests for the first time will create an Observability trace like the image below, where you can see spans for the API call, the queue messaging, the PokeAPI (external API) call and database calls. One interesting thing about this trace is that ",(0,t.jsx)(s.strong,{children:"you can observe the entire use case, end to end"}),":"]}),"\n",(0,t.jsx)(s.p,{children:(0,t.jsx)(s.img,{src:n(92471).A+"",width:"1326",height:"1460"})}),"\n",(0,t.jsx)(s.h3,{id:"assertions",children:"Assertions"}),"\n",(0,t.jsxs)(s.p,{children:["With this trace, we can build ",(0,t.jsx)(s.a,{href:"/concepts/assertions",children:"assertions"})," on Tracetest and validate the API and Worker behaviors:"]}),"\n",(0,t.jsxs)(s.ul,{children:["\n",(0,t.jsxs)(s.li,{children:["\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.strong,{children:"The API should enqueue an import task and return HTTP 200 OK:"}),"\n",(0,t.jsx)(s.img,{src:n(68379).A+"",width:"2976",height:"762"}),"\n",(0,t.jsx)(s.img,{src:n(13263).A+"",width:"2982",height:"796"})]}),"\n"]}),"\n",(0,t.jsxs)(s.li,{children:["\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.strong,{children:"The worker should dequeue the import task:"}),"\n",(0,t.jsx)(s.img,{src:n(98035).A+"",width:"2958",height:"882"})]}),"\n"]}),"\n",(0,t.jsxs)(s.li,{children:["\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.strong,{children:"PokeAPI should return a valid response:"}),"\n",(0,t.jsx)(s.img,{src:n(38507).A+"",width:"2976",height:"828"})]}),"\n"]}),"\n",(0,t.jsxs)(s.li,{children:["\n",(0,t.jsxs)(s.p,{children:[(0,t.jsx)(s.strong,{children:"The database should respond with low latency (< 200ms):"}),"\n",(0,t.jsx)(s.img,{src:n(23974).A+"",width:"2972",height:"790"})]}),"\n"]}),"\n"]}),"\n",(0,t.jsx)(s.p,{children:"Now you can validate this entire use case."}),"\n",(0,t.jsx)(s.h3,{id:"test-definition",children:"Test Definition"}),"\n",(0,t.jsxs)(s.p,{children:["If you want to replicate this entire test on Tracetest, you can replicate these steps on our Web UI or using our CLI, saving the following test definition as the file ",(0,t.jsx)(s.code,{children:"test-definition.yml"})," and later running:"]}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-sh",children:"tracetest run test -f test-definition.yml\n"})}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-yaml",children:'type: Test\nspec:\n  name: Pokeshop - Import\n  description: Import a Pokemon\n  trigger:\n    type: http\n    httpRequest:\n      url: http://demo-pokemon-api.demo/pokemon/import\n      method: POST\n      headers:\n      - key: Content-Type\n        value: application/json\n      body: \'{"id":52}\'\n  specs:\n  - selector: span[tracetest.span.type="messaging" name="queue.synchronizePokemon\n      send" messaging.system="rabbitmq" messaging.destination="queue.synchronizePokemon"]\n    assertions:\n    - attr:messaging.payload = \'{"id":52}\'\n  - selector: span[tracetest.span.type="http" name="POST /pokemon/import" http.method="POST"]\n    assertions:\n    - attr:http.status_code = 200\n    - attr:http.response.body = \'{"id":52}\'\n  - selector: span[tracetest.span.type="messaging" name="queue.synchronizePokemon\n      receive" messaging.system="rabbitmq" messaging.destination="queue.synchronizePokemon"]\n    assertions:\n    - attr:name = "queue.synchronizePokemon receive"\n  - selector: span[tracetest.span.type="http" name="HTTP GET pokeapi.pokemon" http.method="GET"]\n    assertions:\n    - attr:http.response.body  =  \'{"name":"meowth"}\'\n    - attr:http.status_code  =  200\n  - selector: span[tracetest.span.type="database"]\n    assertions:\n    - attr:tracetest.span.duration <= 200ms\n\n'})})]})}function d(e={}){const{wrapper:s}={...(0,o.R)(),...e.components};return s?(0,t.jsx)(s,{...e,children:(0,t.jsx)(l,{...e})}):l(e)}},13263:(e,s,n)=>{n.d(s,{A:()=>t});const t=n.p+"assets/images/import-pokemon-api-test-spec-9d675cd08e1e5adb7ee0f4db28fd604d.png"},23974:(e,s,n)=>{n.d(s,{A:()=>t});const t=n.p+"assets/images/import-pokemon-db-latency-test-spec-e52791c95be226fa0e1402c4486be532.png"},98035:(e,s,n)=>{n.d(s,{A:()=>t});const t=n.p+"assets/images/import-pokemon-message-dequeue-test-spec-f8a3e985984c2d532dbf86023eef120d.png"},68379:(e,s,n)=>{n.d(s,{A:()=>t});const t=n.p+"assets/images/import-pokemon-message-enqueue-test-spec-32297657baf2bd7c228969234177407f.png"},38507:(e,s,n)=>{n.d(s,{A:()=>t});const t=n.p+"assets/images/import-pokemon-pokeapi-call-test-spec-383afa38e11a20f031a47330f9a8ce38.png"},92471:(e,s,n)=>{n.d(s,{A:()=>t});const t=n.p+"assets/images/import-pokemon-trace-f9d7ab3a3bbcba45ff785aab36fab085.png"},28453:(e,s,n)=>{n.d(s,{R:()=>i,x:()=>r});var t=n(96540);const o={},a=t.createContext(o);function i(e){const s=t.useContext(a);return t.useMemo((function(){return"function"==typeof e?e(s):{...s,...e}}),[s,e])}function r(e){let s;return s=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:i(e.components),t.createElement(a.Provider,{value:s},e.children)}}}]);