"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[797],{10007:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>l,contentTitle:()=>i,default:()=>h,frontMatter:()=>r,metadata:()=>o,toc:()=>c});var a=n(74848),s=n(28453);const r={id:"artillery-engine",title:"Performance and Trace-Based Tests with Tracetest and Artillery Engine",description:"Quickstart on how to use the Tracetest x Artillery Integration to enhance Performance Tests with Trace-Based Testing using Tracetest.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","artillery","load testing","performance testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},i=void 0,o={id:"tools-and-integrations/artillery-engine",title:"Performance and Trace-Based Tests with Tracetest and Artillery Engine",description:"Quickstart on how to use the Tracetest x Artillery Integration to enhance Performance Tests with Trace-Based Testing using Tracetest.",source:"@site/docs/tools-and-integrations/artillery-engine.mdx",sourceDirName:"tools-and-integrations",slug:"/tools-and-integrations/artillery-engine",permalink:"/tools-and-integrations/artillery-engine",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/tools-and-integrations/artillery-engine.mdx",tags:[],version:"current",frontMatter:{id:"artillery-engine",title:"Performance and Trace-Based Tests with Tracetest and Artillery Engine",description:"Quickstart on how to use the Tracetest x Artillery Integration to enhance Performance Tests with Trace-Based Testing using Tracetest.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","artillery","load testing","performance testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"guidesSidebar",previous:{title:"Performance and Trace-Based Tests with Tracetest and Artillery Plugin",permalink:"/tools-and-integrations/artillery-plugin"},next:{title:"Performance and Trace-Based Tests with Tracetest and k6",permalink:"/tools-and-integrations/k6"}},l={},c=[{value:"Why is this important?",id:"why-is-this-important",level:2},{value:"The Tracetest Artillery Integration NPM Packages",id:"the-tracetest-artillery-integration-npm-packages",level:2},{value:"Today You&#39;ll Learn How to integrate Trace-Based Tests with your Aritllery Test Scripts",id:"today-youll-learn-how-to-integrate-trace-based-tests-with-your-aritllery-test-scripts",level:2},{value:"Requirements",id:"requirements",level:2},{value:"Run This Quckstart Example",id:"run-this-quckstart-example",level:2},{value:"Project Structure",id:"project-structure",level:2},{value:"Installing the <code>artillery-engine-tracetest</code> NPM Package",id:"installing-the-artillery-engine-tracetest-npm-package",level:2},{value:"Tracetest Test Definitions",id:"tracetest-test-definitions",level:2},{value:"Creating the Artillery Test Script",id:"creating-the-artillery-test-script",level:2},{value:"Running the Full Example",id:"running-the-full-example",level:2},{value:"Finding the Results",id:"finding-the-results",level:2},{value:"What&#39;s Next?",id:"whats-next",level:2},{value:"Learn More",id:"learn-more",level:2}];function p(e){const t={a:"a",admonition:"admonition",code:"code",h2:"h2",img:"img",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,s.R)(),...e.components};return(0,a.jsxs)(a.Fragment,{children:[(0,a.jsx)(t.admonition,{title:"Commercial Feature",type:"tip",children:(0,a.jsx)(t.p,{children:(0,a.jsx)(t.a,{href:"https://tracetest.io/pricing",children:"Feature available only in Cloud-based Managed Tracetest & Enterprise Self-hosted Tracetest."})})}),"\n",(0,a.jsx)(t.admonition,{title:"Version Compatibility",type:"info",children:(0,a.jsxs)(t.p,{children:["This integration is compatible with ",(0,a.jsx)(t.a,{href:"https://github.com/artilleryio/artillery/releases/tag/artillery-2.0.10",children:"Artillery v2.0.10"})," and above."]})}),"\n",(0,a.jsx)(t.admonition,{type:"note",children:(0,a.jsx)(t.p,{children:(0,a.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples/quick-start-artillery",children:"Check out the source code on GitHub here."})})}),"\n",(0,a.jsxs)(t.p,{children:[(0,a.jsx)(t.a,{href:"https://app.tracetest.io/",children:"Tracetest"})," is a testing tool based on ",(0,a.jsx)(t.a,{href:"https://opentelemetry.io/",children:"OpenTelemetry"})," that allows you to test your distributed application. It allows you to use data from distributed traces generated by OpenTelemetry to validate and assert if your application has the desired behavior defined by your test definitions."]}),"\n",(0,a.jsxs)(t.p,{children:[(0,a.jsx)(t.a,{href:"https://artillery.io/",children:"Artillery"})," is a modern, powerful load-testing toolkit. Artillery is designed to help developers and testers simulate traffic to their applications, APIs, and microservices. It allows users to define scenarios to test how their systems behave under different loads."]}),"\n",(0,a.jsx)(t.h2,{id:"why-is-this-important",children:"Why is this important?"}),"\n",(0,a.jsx)(t.p,{children:"Artillery is it's a great tool in its own right that allows you to replicate most of the production challenges you might encounter. But, as with all of the tools that only test the initial transaction between the client side and the server, you can only run validations against the immediate response from the service."}),"\n",(0,a.jsx)(t.h2,{id:"the-tracetest-artillery-integration-npm-packages",children:"The Tracetest Artillery Integration NPM Packages"}),"\n",(0,a.jsxs)(t.p,{children:["The ",(0,a.jsx)(t.a,{href:"https://www.npmjs.com/package/artillery-engine-tracetest",children:(0,a.jsx)(t.code,{children:"artillery-engine-tracetest"})})," package enables you to run performance tests using Tracetest's as the main engine. Reaching directly to your APIs and fetching the generated telemetry data to be used to run trace-based tests."]}),"\n",(0,a.jsx)(t.admonition,{type:"note",children:(0,a.jsxs)(t.p,{children:["Are you using the HTTP engine for your Artillery Test Scripts? Take a look at the ",(0,a.jsx)(t.a,{href:"./artillery-plugin",children:(0,a.jsx)(t.code,{children:"artillery-plugin-tracetest"})})," guide on how to seemingly integrate it with your existing setup."]})}),"\n",(0,a.jsx)(t.h2,{id:"today-youll-learn-how-to-integrate-trace-based-tests-with-your-aritllery-test-scripts",children:"Today You'll Learn How to integrate Trace-Based Tests with your Aritllery Test Scripts"}),"\n",(0,a.jsxs)(t.p,{children:["This is a simple quick-start guide on how to use the Tracetest ",(0,a.jsx)(t.code,{children:"artillery-engine-tracetest"})," NPM package to enhance your Artillery Test Scripts with trace-based testing. The infrastructure will use the Pokeshop Demo as a testing ground, triggering requests against it and generating telemetry data."]}),"\n",(0,a.jsx)(t.h2,{id:"requirements",children:"Requirements"}),"\n",(0,a.jsxs)(t.p,{children:[(0,a.jsx)(t.strong,{children:"Tracetest Account"}),":"]}),"\n",(0,a.jsxs)(t.ul,{children:["\n",(0,a.jsxs)(t.li,{children:["Sign up to ",(0,a.jsx)(t.a,{href:"https://app.tracetest.io",children:(0,a.jsx)(t.code,{children:"app.tracetest.io"})})," or follow the ",(0,a.jsx)(t.a,{href:"/getting-started/overview",children:"get started"})," docs."]}),"\n",(0,a.jsxs)(t.li,{children:["Create an ",(0,a.jsx)(t.a,{href:"/concepts/environments",children:"environment"}),"."]}),"\n",(0,a.jsxs)(t.li,{children:["Create an ",(0,a.jsx)(t.a,{href:"/concepts/environment-tokens",children:"environment token"}),"."]}),"\n",(0,a.jsxs)(t.li,{children:["Have access to the environment's ",(0,a.jsx)(t.a,{href:"/configuration/agent",children:"agent API key"}),"."]}),"\n"]}),"\n",(0,a.jsxs)(t.p,{children:[(0,a.jsx)(t.strong,{children:"Docker"}),": Have ",(0,a.jsx)(t.a,{href:"https://docs.docker.com/get-docker/",children:"Docker"})," and ",(0,a.jsx)(t.a,{href:"https://docs.docker.com/compose/install/",children:"Docker Compose"})," installed on your machine."]}),"\n",(0,a.jsx)(t.h2,{id:"run-this-quckstart-example",children:"Run This Quckstart Example"}),"\n",(0,a.jsx)(t.p,{children:"The example below is provided as part of the Tracetest project. You can download and run the example by following these steps:"}),"\n",(0,a.jsx)(t.p,{children:"Clone the Tracetest project and go to the TypeScript Quickstart:"}),"\n",(0,a.jsx)(t.pre,{children:(0,a.jsx)(t.code,{className:"language-bash",children:"git clone https://github.com/kubeshop/tracetest\ncd tracetest/examples/quick-start-artillery\n"})}),"\n",(0,a.jsx)(t.p,{children:"Follow these instructions to run the included demo app and TypeScript example:"}),"\n",(0,a.jsxs)(t.ol,{children:["\n",(0,a.jsxs)(t.li,{children:["Copy the ",(0,a.jsx)(t.code,{children:".env.template"})," file to ",(0,a.jsx)(t.code,{children:".env"}),"."]}),"\n",(0,a.jsxs)(t.li,{children:["Log into the ",(0,a.jsx)(t.a,{href:"https://app.tracetest.io/",children:"Tracetest app"}),"."]}),"\n",(0,a.jsxs)(t.li,{children:["This example is configured to use Jaeger. Ensure the environment you will be utilizing to run this example is also configured to use the Jaeger Tracing Backend by clicking on Settings, Tracing Backend, Jaeger, updating the url to ",(0,a.jsx)(t.code,{children:"jaeger:16685"}),", Test Connection and Save."]}),"\n",(0,a.jsxs)(t.li,{children:["Fill out the ",(0,a.jsx)(t.a,{href:"https://docs.tracetest.io/concepts/environment-tokens",children:"token"})," and ",(0,a.jsx)(t.a,{href:"https://docs.tracetest.io/concepts/agent",children:"agent API key"})," details by editing your ",(0,a.jsx)(t.code,{children:".env"})," file. You can find these values in the Settings area for your environment."]}),"\n",(0,a.jsxs)(t.li,{children:["Run ",(0,a.jsx)(t.code,{children:"docker compose up -d"}),"."]}),"\n",(0,a.jsxs)(t.li,{children:["Run ",(0,a.jsx)(t.code,{children:"npm i"})," to install the required dependencies."]}),"\n",(0,a.jsxs)(t.li,{children:["Run ",(0,a.jsx)(t.code,{children:"npm run test:engine"})," to run the example."]}),"\n",(0,a.jsx)(t.li,{children:"The output will show the test results and the Tracetest URL for each test run."}),"\n",(0,a.jsx)(t.li,{children:"Follow the links in the log to view the test runs programmatically created by the Atillery execution."}),"\n"]}),"\n",(0,a.jsx)(t.p,{children:"Follow along with the sections below for an in detail breakdown of what the example you just ran did and how it works."}),"\n",(0,a.jsx)(t.h2,{id:"project-structure",children:"Project Structure"}),"\n",(0,a.jsx)(t.p,{children:"The quick-start Artillery project is built with Docker Compose."}),"\n",(0,a.jsxs)(t.p,{children:["The ",(0,a.jsx)(t.a,{href:"/live-examples/pokeshop/overview",children:"Pokeshop Demo App"})," is a complete example of a distributed application using different back-end and front-end services. We will be launching it and running tests against it as part of this example."]}),"\n",(0,a.jsxs)(t.p,{children:["The ",(0,a.jsx)(t.code,{children:"docker-compose.yaml"})," file in the root directory of the quick start runs the Pokeshop Demo app, the OpenTelemetry Collector, Jaeger, and the ",(0,a.jsx)(t.a,{href:"/concepts/agent",children:"Tracetest Agent"})," setup."]}),"\n",(0,a.jsxs)(t.p,{children:["The Artillery Plugin quick start has two primary files: a Test Script file ",(0,a.jsx)(t.code,{children:"engine-test.yaml"})," that defines the Artillery execution, and a Tracetest Definition file ",(0,a.jsx)(t.code,{children:"import-pokemon.yaml"})," that contains the specs and execution of the trace-based tests."]}),"\n",(0,a.jsxs)(t.h2,{id:"installing-the-artillery-engine-tracetest-npm-package",children:["Installing the ",(0,a.jsx)(t.code,{children:"artillery-engine-tracetest"})," NPM Package"]}),"\n",(0,a.jsxs)(t.p,{children:["The first step when using the Artillery Plugin NPM package is to install the ",(0,a.jsx)(t.code,{children:"artillery-engine-tracetest"})," NPM Package. It is as easy as running the following command:"]}),"\n",(0,a.jsx)(t.pre,{children:(0,a.jsx)(t.code,{className:"language-bash",children:"npm i artillery-engine-tracetest\n"})}),"\n",(0,a.jsxs)(t.p,{children:["Once you have installed the ",(0,a.jsx)(t.code,{children:"artillery-engine-tracetest"})," package, you can use it as part of your Artillery Test Scripts to trigger trace-based tests and run checks against the resulting telemetry data."]}),"\n",(0,a.jsx)(t.h2,{id:"tracetest-test-definitions",children:"Tracetest Test Definitions"}),"\n",(0,a.jsxs)(t.p,{children:["The ",(0,a.jsx)(t.code,{children:"import-pokemon.yaml"})," file contains the YAML version of the test definitions that will be used to run the tests. It uses the Artillery trigger to execute requests against the Pokeshop Demo."]}),"\n",(0,a.jsx)(t.pre,{children:(0,a.jsx)(t.code,{className:"language-yaml",metastring:'title="import-pokemon.yaml"',children:'type: Test\nspec:\n  id: artillery-engine-import-pokemon\n  name: "Artillery Engine: Import Pokemon"\n\n  # Executing a POST request to the /pokemon/import endpoint\n  trigger:\n    type: http\n    httpRequest:\n      method: POST\n      url: ${var:ENDPOINT}/pokemon/import\n      body: \'{"id": ${var:POKEMON_ID}}\'\n      headers:\n        - key: Content-Type\n          value: application/json\n\n  # Defining the trace-based tests\n  specs:\n    - selector: span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"]\n      name: "All HTTP Spans: Status  code is 200"\n      assertions:\n        - attr:http.status_code = 200\n    - selector: span[tracetest.span.type="http" name="GET" http.method="GET"]\n      assertions:\n        - attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n    - selector: span[tracetest.span.type="database"]\n      name: "All Database Spans: Processing time is less than 1s"\n      assertions:\n        - attr:tracetest.span.duration < 1s\n\n  # Defining the outputs\n  outputs:\n    - name: DATABASE_POKEMON_ID\n      selector: span[tracetest.span.type="database" name="create postgres.pokemon" db.system="postgres" db.name="postgres" db.user="postgres" db.operation="create" db.sql.table="pokemon"]\n      value: attr:db.result | json_path \'$.id\'\n'})}),"\n",(0,a.jsx)(t.h2,{id:"creating-the-artillery-test-script",children:"Creating the Artillery Test Script"}),"\n",(0,a.jsxs)(t.p,{children:["The ",(0,a.jsx)(t.code,{children:"engine-test.yaml"})," file contains the Artillery Test Script that will be used to trigger requests against the Pokeshop Demo and run trace-based tests. The steps executed by this script are the following:"]}),"\n",(0,a.jsxs)(t.ol,{children:["\n",(0,a.jsxs)(t.li,{children:["Congures a phase that will execute ",(0,a.jsx)(t.code,{children:"10"})," Tracetest test runs."]}),"\n",(0,a.jsxs)(t.li,{children:["Includes the API ",(0,a.jsx)(t.code,{children:"token"})," to access the Tracetest APIs."]}),"\n",(0,a.jsxs)(t.li,{children:["Defines the scenarios using the ",(0,a.jsx)(t.code,{children:"tracetest"})," engine adding the Test definition from the ",(0,a.jsx)(t.code,{children:"import-pokemon.yaml"})," which is an HTTP request to the ",(0,a.jsx)(t.code,{children:"POST pokemon/import"})," endpoint sending ",(0,a.jsx)(t.code,{children:"6"})," (Charizard) as the Pokemon ID."]}),"\n",(0,a.jsxs)(t.li,{children:["The summary format is set to ",(0,a.jsx)(t.code,{children:"pretty"})," to display the results in the console."]}),"\n"]}),"\n",(0,a.jsx)(t.pre,{children:(0,a.jsx)(t.code,{className:"language-yaml",metastring:'title="engine-test.yaml"',children:'config:\n  target: my_target\n  tracetest:\n    token: <your-api-token>\n  phases:\n    - duration: 2\n      arrivalRate: 5\n  engines:\n    tracetest: {}\nscenarios:\n  - name: tracetest_engine_test\n    engine: tracetest\n    flow:\n      - test:\n          definition: import-pokemon.yaml\n          runInfo:\n            variables:\n              - key: ENDPOINT\n                value: http://api:8081\n              - key: POKEMON_ID\n                value: "6"\n      - summary:\n          format: "pretty"\n'})}),"\n",(0,a.jsx)(t.h2,{id:"running-the-full-example",children:"Running the Full Example"}),"\n",(0,a.jsx)(t.p,{children:"To start the full setup, run the following command:"}),"\n",(0,a.jsx)(t.pre,{children:(0,a.jsx)(t.code,{className:"language-bash",children:"docker compose up -d\nnpm run test:engine\n"})}),"\n",(0,a.jsx)(t.h2,{id:"finding-the-results",children:"Finding the Results"}),"\n",(0,a.jsx)(t.p,{children:"The output from the Tracetest Engine script should be visible in the console log after running the test command. This log will show links to Tracetest for each of the test runs invoked by the Artillery Testing Script. Click a link to launch Tracetest and view the test result."}),"\n",(0,a.jsx)(t.pre,{children:(0,a.jsx)(t.code,{className:"language-bash",children:'\n> quick-start-artillery@1.0.0 test:engine\n> artillery run engine-test.yaml\n\nTest run id: t47ez_f4gmj899tr4nnp465h38d4mnker7z_w8yp\nPhase started: unnamed (index: 0, duration: 2s) 09:29:33(-0600)\n\nPhase completed: unnamed (index: 0, duration: 2s) 09:29:35(-0600)\n\n--------------------------------------\nMetrics for period to: 09:29:40(-0600) (width: 1.203s)\n--------------------------------------\n\nvusers.created: ................................................................ 10\nvusers.created_by_name.tracetest_engine_test: .................................. 10\n\n\nWarning: multiple batches of metrics for period 1709825370000 2024-03-07T15:29:30.000Z\n\u280f \u2718 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/70) - trace id: 4a2e41bf5b497a643056c8a08d122d82\n  > All HTTP Spans: Status  code is 200 (span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"])\n\n    \u2022 Expected: attr:http.status_code = 200\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.status_code" from. To fix this error, create a selector matching at least one span.\n  > span[tracetest.span.type="http" name="GET" http.method="GET"]\n\n    \u2022 Expected: attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.route" from. To fix this error, create a selector matching at least one span.\n  > All Database Spans: Processing time is less than 1s (span[tracetest.span.type="database"])\n\u2819 \u2718 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/74) - trace id: a00ff614a06be196ec5cf74ace906b3e\n  > All HTTP Spans: Status  code is 200 (span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"])\n\n    \u2022 Expected: attr:http.status_code = 200\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.status_code" from. To fix this error, create a selector matching at least one span.\n  > span[tracetest.span.type="http" name="GET" http.method="GET"]\n\n    \u2022 Expected: attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.route" from. To fix this error, create a selector matching at least one span.\n  > All Database Spans: Processing time is less than 1s (span[tracetest.span.type="database"])\n\u2838 \u2718 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/76) - trace id: 1fb3ded60f582ecfd151a019df2479cd\n  > All HTTP Spans: Status  code is 200 (span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"])\n\n    \u2022 Expected: attr:http.status_code = 200\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.status_code" from. To fix this error, create a selector matching at least one span.\n  > span[tracetest.span.type="http" name="GET" http.method="GET"]\n\n    \u2022 Expected: attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.route" from. To fix this error, create a selector matching at least one span.\n  > All Database Spans: Processing time is less than 1s (span[tracetest.span.type="database"])\n\u2839 \u2718 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/78) - trace id: 750df8c1d55ed6a4e18f4b2b8b573a14\n  > All HTTP Spans: Status  code is 200 (span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"])\n\n    \u2022 Expected: attr:http.status_code = 200\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.status_code" from. To fix this error, create a selector matching at least one span.\n  > span[tracetest.span.type="http" name="GET" http.method="GET"]\n\n    \u2022 Expected: attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.route" from. To fix this error, create a selector matching at least one span.\n  > All Database Spans: Processing time is less than 1s (span[tracetest.span.type="database"])\n\u2826 \u2714 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/71) - trace id: c3bdcfb86ca77a9455b29882f593848a\n\u280f \u2718 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/79) - trace id: 76ecf991c6eb37c2a0dee857a699dc76\n  > All HTTP Spans: Status  code is 200 (span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"])\n\n    \u2022 Expected: attr:http.status_code = 200\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.status_code" from. To fix this error, create a selector matching at least one span.\n  > span[tracetest.span.type="http" name="GET" http.method="GET"]\n\n    \u2022 Expected: attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.route" from. To fix this error, create a selector matching at least one span.\n  > All Database Spans: Processing time is less than 1s (span[tracetest.span.type="database"])\n\u2819 \u2718 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/72) - trace id: cece9c79434b3313bf23150bc967e82f\n  > All HTTP Spans: Status  code is 200 (span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"])\n\n    \u2022 Expected: attr:http.status_code = 200\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.status_code" from. To fix this error, create a selector matching at least one span.\n  > span[tracetest.span.type="http" name="GET" http.method="GET"]\n\n    \u2022 Expected: attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.route" from. To fix this error, create a selector matching at least one span.\n  > All Database Spans: Processing time is less than 1s (span[tracetest.span.type="database"])\n\u2807 \u2718 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/77) - trace id: b3f3682eb328195ede902177ff4170a1\n  > All HTTP Spans: Status  code is 200 (span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"])\n\n    \u2022 Expected: attr:http.status_code = 200\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.status_code" from. To fix this error, create a selector matching at least one span.\n  > span[tracetest.span.type="http" name="GET" http.method="GET"]\n\n    \u2022 Expected: attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.route" from. To fix this error, create a selector matching at least one span.\n  > All Database Spans: Processing time is less than 1s (span[tracetest.span.type="database"])\n\u2838 \u2718 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/75) - trace id: c2d4624f20060565b0b5e9f97ad7e7a7\n  > All HTTP Spans: Status  code is 200 (span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"])\n\n    \u2022 Expected: attr:http.status_code = 200\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.status_code" from. To fix this error, create a selector matching at least one span.\n  > span[tracetest.span.type="http" name="GET" http.method="GET"]\n\n    \u2022 Expected: attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.route" from. To fix this error, create a selector matching at least one span.\n  > All Database Spans: Processing time is less than 1s (span[tracetest.span.type="database"])\n\u2819 \u2718 Artillery Engine: Import Pokemon (https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/environments/ttenv_231b49e808c29e6a/test/artillery-plugin-import-pokemon/run/73) - trace id: 878bb73b5234f764f4d289bf657b40cf\n  > All HTTP Spans: Status  code is 200 (span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"])\n\n    \u2022 Expected: attr:http.status_code = 200\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.status_code" from. To fix this error, create a selector matching at least one span.\n  > span[tracetest.span.type="http" name="GET" http.method="GET"]\n\n    \u2022 Expected: attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\n\n      No Spans\n      Received: resolution error: there are no matching spans to retrieve the attribute "http.route" from. To fix this error, create a selector matching at least one span.\n  > All Database Spans: Processing time is less than 1s (span[tracetest.span.type="database"])\n--------------------------------------\nMetrics for period to: 09:30:00(-0600) (width: 0.015s)\n--------------------------------------\n\ntracetest.tests_failed: ........................................................ 1\nvusers.completed: .............................................................. 1\nvusers.failed: ................................................................. 0\nvusers.session_length:\n  min: ......................................................................... 21797.1\n  max: ......................................................................... 21797.1\n  mean: ........................................................................ 21797.1\n  median: ...................................................................... 21813.5\n  p95: ......................................................................... 21813.5\n  p99: ......................................................................... 21813.5\n\n\nAll VUs finished. Total time: 34 seconds\n\n--------------------------------\nSummary report @ 09:30:09(-0600)\n--------------------------------\n\ntracetest.tests_failed: ........................................................ 9\ntracetest.tests_succeeded: ..................................................... 1\nvusers.completed: .............................................................. 10\nvusers.created: ................................................................ 10\nvusers.created_by_name.tracetest_engine_test: .................................. 10\nvusers.failed: ................................................................. 0\nvusers.session_length:\n  min: ......................................................................... 21797.1\n  max: ......................................................................... 33447\n  mean: ........................................................................ 28651.1\n  median: ...................................................................... 28862.3\n  p95: ......................................................................... 32542.3\n  p99: ......................................................................... 32542.3\n'})}),"\n",(0,a.jsx)(t.admonition,{type:"note",children:(0,a.jsx)(t.p,{children:"Most of the tests will fail as the import Pokemon flow reads from memory if the info already exists. This is expected behavior."})}),"\n",(0,a.jsx)(t.admonition,{title:"View these tests in our Demo environment",type:"tip",children:(0,a.jsx)(t.p,{children:(0,a.jsx)(t.a,{href:"https://app.tracetest.io/organizations/ttorg_2179a9cd8ba8dfa5/invites/invite_760904a64b4b9dc9/accept",children:"\ud83d\udc49 Join our shared Pokeshop API Demo environment."})})}),"\n",(0,a.jsx)(t.h2,{id:"whats-next",children:"What's Next?"}),"\n",(0,a.jsx)(t.p,{children:"After running the initial set of tests, you can click the run link for any of them, update the assertions, and run the scripts once more. This flow enables complete a trace-based TDD flow."}),"\n",(0,a.jsx)(t.p,{children:(0,a.jsx)(t.img,{alt:"assertions",src:n(31681).A+"",width:"960",height:"553"})}),"\n",(0,a.jsx)(t.h2,{id:"learn-more",children:"Learn More"}),"\n",(0,a.jsxs)(t.p,{children:["Please visit our ",(0,a.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples",children:"examples in GitHub"})," and join our ",(0,a.jsx)(t.a,{href:"https://dub.sh/tracetest-community",children:"Slack Community"})," for more info!"]})]})}function h(e={}){const{wrapper:t}={...(0,s.R)(),...e.components};return t?(0,a.jsx)(t,{...e,children:(0,a.jsx)(p,{...e})}):p(e)}},31681:(e,t,n)=>{n.d(t,{A:()=>a});const a=n.p+"assets/images/tracetest-cloud-typescript-resize-bd6c9a72c3825d10494de1a1fbf3171c.gif"},28453:(e,t,n)=>{n.d(t,{R:()=>i,x:()=>o});var a=n(96540);const s={},r=a.createContext(s);function i(e){const t=a.useContext(r);return a.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function o(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:i(e.components),a.createElement(r.Provider,{value:t},e.children)}}}]);