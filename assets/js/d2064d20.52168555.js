"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[7013],{37628:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>o,default:()=>p,frontMatter:()=>i,metadata:()=>a,toc:()=>l});var s=n(74848),r=n(28453);const i={id:"running-tracetest-with-azure-app-insights-pokeshop",title:"Pokeshop API and Azure Application Insights with OpenTelemetry Collector",description:"Quick start on how to configure the Pokeshop API Demo with OpenTelemetry traces, Azure Application Insights as a trace data store, including the OpenTelemetry Collector, and Tracetest for enhancing your E2E and integration tests with trace-based testing.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","azure functions","azure app insights","azure application insights","azure tracing","azure monitor","opentelemetry"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},o=void 0,a={id:"examples-tutorials/recipes/running-tracetest-with-azure-app-insights-pokeshop",title:"Pokeshop API and Azure Application Insights with OpenTelemetry Collector",description:"Quick start on how to configure the Pokeshop API Demo with OpenTelemetry traces, Azure Application Insights as a trace data store, including the OpenTelemetry Collector, and Tracetest for enhancing your E2E and integration tests with trace-based testing.",source:"@site/docs/examples-tutorials/recipes/running-tracetest-with-azure-app-insights-pokeshop.mdx",sourceDirName:"examples-tutorials/recipes",slug:"/examples-tutorials/recipes/running-tracetest-with-azure-app-insights-pokeshop",permalink:"/examples-tutorials/recipes/running-tracetest-with-azure-app-insights-pokeshop",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/examples-tutorials/recipes/running-tracetest-with-azure-app-insights-pokeshop.mdx",tags:[],version:"current",frontMatter:{id:"running-tracetest-with-azure-app-insights-pokeshop",title:"Pokeshop API and Azure Application Insights with OpenTelemetry Collector",description:"Quick start on how to configure the Pokeshop API Demo with OpenTelemetry traces, Azure Application Insights as a trace data store, including the OpenTelemetry Collector, and Tracetest for enhancing your E2E and integration tests with trace-based testing.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","azure functions","azure app insights","azure application insights","azure tracing","azure monitor","opentelemetry"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"guidesSidebar",previous:{title:"Node.js and Azure Application Insights with OpenTelemetry Collector",permalink:"/examples-tutorials/recipes/running-tracetest-with-azure-app-insights-collector"},next:{title:"Node.js, Elasticsearch and Elastic APM",permalink:"/examples-tutorials/recipes/running-tracetest-with-elasticapm"}},c={},l=[{value:"Pokeshop API with Azure Application Insights, OpenTelemetry Collector and Tracetest",id:"pokeshop-api-with-azure-application-insights-opentelemetry-collector-and-tracetest",level:2},{value:"Prerequisites",id:"prerequisites",level:2},{value:"Run This Quckstart Example",id:"run-this-quckstart-example",level:2},{value:"Pokeshop API",id:"pokeshop-api",level:2},{value:"Configuring the Pokeshop Demo API",id:"configuring-the-pokeshop-demo-api",level:2},{value:"Running the Pokeshop Demo API, OpenTelemetry Collector, and Tracetest",id:"running-the-pokeshop-demo-api-opentelemetry-collector-and-tracetest",level:2},{value:"The Test File",id:"the-test-file",level:3},{value:"Learn More",id:"learn-more",level:2}];function h(e){const t={a:"a",admonition:"admonition",code:"code",h2:"h2",h3:"h3",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,r.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(t.admonition,{type:"note",children:(0,s.jsx)(t.p,{children:(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-azure-app-insights-pokeshop",children:"Check out the source code on GitHub here."})})}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://tracetest.io/",children:"Tracetest"})," is a testing tool based on ",(0,s.jsx)(t.a,{href:"https://opentelemetry.io/",children:"OpenTelemetry"})," that allows you to test your distributed application. It allows you to use data from distributed traces generated by OpenTelemetry to validate and assert if your application has the desired behavior defined by your test definitions."]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://learn.microsoft.com/en-us/azure/azure-monitor/app/app-insights-overview",children:"Azure Application Insights"})," is an extension of Azure Monitor and provides application performance monitoring (APM) features. APM tools are useful to monitor applications from development, through test, and into production in the following ways:"]}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsx)(t.li,{children:"Proactively understand how an application is performing."}),"\n",(0,s.jsx)(t.li,{children:"Reactively review application execution data to determine the cause of an incident."}),"\n"]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://github.com/open-telemetry/opentelemetry-collector-contrib",children:"OpenTelemetry Collector Contrib"})," - The official OpenTelemetry Distribution for packages outside of the core collector."]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://docs.tracetest.io/live-examples/pokeshop/overview",children:"Pokeshop API"})," As a testing ground, the team at Tracetest has implemented a sample instrumented API around the ",(0,s.jsx)(t.a,{href:"https://pokeapi.co/",children:"PokeAPI"}),"."]}),"\n",(0,s.jsx)(t.h2,{id:"pokeshop-api-with-azure-application-insights-opentelemetry-collector-and-tracetest",children:"Pokeshop API with Azure Application Insights, OpenTelemetry Collector and Tracetest"}),"\n",(0,s.jsx)(t.p,{children:"This is a simple quick start guide on how to configure a fully instrumented API to be used with Tracetest for enhancing your E2E and integration tests with trace-based testing. The infrastructure will use Azure App Insights as the trace data store, the OpenTelemetry Collector as a middleware and the Pokeshop API to generate the telemetry data."}),"\n",(0,s.jsx)(t.h2,{id:"prerequisites",children:"Prerequisites"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"Tracetest Account"}),":"]}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsxs)(t.li,{children:["Sign up to ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io",children:(0,s.jsx)(t.code,{children:"app.tracetest.io"})})," or follow the ",(0,s.jsx)(t.a,{href:"/getting-started/overview",children:"get started"})," docs."]}),"\n",(0,s.jsxs)(t.li,{children:["Have access to the environment's ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io/retrieve-token",children:"agent API key"}),"."]}),"\n"]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"Azure Account"}),":"]}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsxs)(t.li,{children:["Sign up to ",(0,s.jsx)(t.a,{href:"https://azure.microsoft.com/en-us",children:"Azure"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Install ",(0,s.jsx)(t.a,{href:"https://learn.microsoft.com/en-us/cli/azure/install-azure-cli",children:"Azure CLI"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Create an Application Insights app and get a ",(0,s.jsx)(t.a,{href:"https://learn.microsoft.com/en-us/azure/bot-service/bot-service-resources-app-insights-keys?view=azure-bot-service-4.0",children:"Instrumentation Key"}),"."]}),"\n"]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"Docker"}),": Have ",(0,s.jsx)(t.a,{href:"https://docs.docker.com/get-docker/",children:"Docker"})," and ",(0,s.jsx)(t.a,{href:"https://docs.docker.com/compose/install/",children:"Docker Compose"})," installed on your machine."]}),"\n",(0,s.jsx)(t.h2,{id:"run-this-quckstart-example",children:"Run This Quckstart Example"}),"\n",(0,s.jsx)(t.p,{children:"The example below is provided as part of the Tracetest project. You can download and run the example by following these steps:"}),"\n",(0,s.jsx)(t.p,{children:"Clone the Tracetest project and go to the Azure Node.js Quickstart:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"git clone https://github.com/kubeshop/tracetest\ncd tracetest/examples/tracetest-azure-app-insights-pokeshop\n"})}),"\n",(0,s.jsx)(t.p,{children:"Follow these instructions to run the quick start:"}),"\n",(0,s.jsxs)(t.ol,{children:["\n",(0,s.jsxs)(t.li,{children:["Copy the ",(0,s.jsx)(t.code,{children:".env.template"})," file to ",(0,s.jsx)(t.code,{children:".env"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Fill out the ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io/retrieve-token",children:"TRACETEST_TOKEN and ENVIRONMENT_ID"})," details by editing your ",(0,s.jsx)(t.code,{children:".env"})," file."]}),"\n",(0,s.jsxs)(t.li,{children:["Fill out the ",(0,s.jsx)(t.a,{href:"https://learn.microsoft.com/en-us/azure/azure-monitor/app/sdk-connection-string?tabs=dotnet5#find-your-connection-string",children:"APP_INSIGHTS_ACCESS_TOKEN, APP_INSIGHTS_ARM_ID and APP_INSIGHTS_INSTRUMENTATION_STRING"})," details by editing your ",(0,s.jsx)(t.code,{children:".env"})," file."]}),"\n",(0,s.jsxs)(t.li,{children:["Run ",(0,s.jsx)(t.code,{children:"docker compose run tracetest-run"}),"."]}),"\n",(0,s.jsx)(t.li,{children:"Follow the links in the output to view the test results."}),"\n"]}),"\n",(0,s.jsx)(t.p,{children:"Follow the sections below for a detailed breakdown of what the example you just ran did and how it works."}),"\n",(0,s.jsx)(t.h2,{id:"pokeshop-api",children:"Pokeshop API"}),"\n",(0,s.jsxs)(t.p,{children:["The Pokeshop API is instrumented using the ",(0,s.jsx)(t.a,{href:"https://opentelemetry.io/docs/instrumentation/js/getting-started/nodejs/",children:"OpenTelemetry standards for Node.js"}),", sending the data to the OpenTelemetry collector that will be pushing the telemetry information to the Azure Cloud."]}),"\n",(0,s.jsx)(t.admonition,{title:"Pokeshop Demo API Architecture",type:"note",children:(0,s.jsx)(t.p,{children:(0,s.jsx)(t.a,{href:"/live-examples/pokeshop/overview",children:"View the full system architecture of the Pokeshop Demo API, here."})})}),"\n",(0,s.jsx)(t.h2,{id:"configuring-the-pokeshop-demo-api",children:"Configuring the Pokeshop Demo API"}),"\n",(0,s.jsxs)(t.p,{children:["Configure the ",(0,s.jsx)(t.code,{children:".env"})," like shown below."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:'# Get the required information here: https://app.tracetest.io/retrieve-token\n\nTRACETEST_TOKEN="<YOUR_TRACETEST_TOKEN>"\nTRACETEST_ENVIRONMENT_ID="<YOUR_ENV_ID>"\n\n# Azure\nAPP_INSIGHTS_ACCESS_TOKEN=""\nAPP_INSIGHTS_ARM_ID="/subscriptions/<id>/resourceGroups/app-insights-1/providers/microsoft.insights/components/<name>"\nAPP_INSIGHTS_INSTRUMENTATION_STRING=""\n'})}),"\n",(0,s.jsx)(t.h2,{id:"running-the-pokeshop-demo-api-opentelemetry-collector-and-tracetest",children:"Running the Pokeshop Demo API, OpenTelemetry Collector, and Tracetest"}),"\n",(0,s.jsxs)(t.p,{children:["The ",(0,s.jsxs)(t.a,{href:"https://github.com/kubeshop/tracetest/blob/main/examples/tracetest-azure-app-insights-pokeshop/docker-compose.yaml",children:[(0,s.jsx)(t.code,{children:"docker-compose.yaml"})," file"]})," in the root directory contains the Node.js app, OpenTelemetry Collector and Tracetest Agent."]}),"\n",(0,s.jsx)(t.p,{children:"To run everything including Tracetest tests, run this command:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"docker compose run tracetest-run\n"})}),"\n",(0,s.jsx)(t.p,{children:"This will:"}),"\n",(0,s.jsxs)(t.ol,{children:["\n",(0,s.jsx)(t.li,{children:"Start the Node.js app and send the traces to Azure App Insights."}),"\n",(0,s.jsx)(t.li,{children:"Start the Tracetest Agent."}),"\n",(0,s.jsx)(t.li,{children:"Configure the Azure App Insights tracing backend and create tests in your environment."}),"\n",(0,s.jsx)(t.li,{children:"Run the tests."}),"\n"]}),"\n",(0,s.jsx)(t.h3,{id:"the-test-file",children:"The Test File"}),"\n",(0,s.jsxs)(t.p,{children:["Check out the ",(0,s.jsx)(t.code,{children:"resources/test.yaml"})," file."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-yaml",children:'# resources/test.yaml\ntype: Test\nspec:\n  id: -ao9stJVg\n  name: Pokeshop - Import\n  description: Import a Pokemon\n  pollingProfile: azure\n  trigger:\n    type: http\n    httpRequest:\n      url: http://demo-api:8081/pokemon/import\n      method: POST\n      headers:\n        - key: Content-Type\n          value: application/json\n      body: \'{"id":6}\'\n  specs:\n    - name: Import Pokemon Span Exists\n      selector: span[tracetest.span.type="general" name="import pokemon"]\n      assertions:\n        - attr:tracetest.selected_spans.count = 1\n    - name: Uses Correct PokemonId\n      selector: span[tracetest.span.type="http" name="GET /pokemon/6" http.method="GET"]\n      assertions:\n        - attr:http.url  =  "https://pokeapi.co/api/v2/pokemon/6"\n    - name: Matching db result with the Pokemon Name\n      selector: span[tracetest.span.type="database" name="create pokeshop.pokemon"]:first\n      assertions:\n        - attr:db.result  contains      "charizard"\n'})}),"\n",(0,s.jsx)(t.h2,{id:"learn-more",children:"Learn More"}),"\n",(0,s.jsxs)(t.p,{children:["Please visit our ",(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples",children:"examples in GitHub"})," and join our ",(0,s.jsx)(t.a,{href:"https://dub.sh/tracetest-community",children:"Slack Community"})," for more info!"]})]})}function p(e={}){const{wrapper:t}={...(0,r.R)(),...e.components};return t?(0,s.jsx)(t,{...e,children:(0,s.jsx)(h,{...e})}):h(e)}},28453:(e,t,n)=>{n.d(t,{R:()=>o,x:()=>a});var s=n(96540);const r={},i=s.createContext(r);function o(e){const t=s.useContext(i);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function a(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:o(e.components),s.createElement(i.Provider,{value:t},e.children)}}}]);