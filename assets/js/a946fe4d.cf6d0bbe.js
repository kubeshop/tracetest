"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[3940],{86098:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>i,default:()=>h,frontMatter:()=>r,metadata:()=>o,toc:()=>l});var s=n(74848),a=n(28453);const r={id:"running-tracetest-without-a-trace-data-store-with-manual-instrumentation",title:"Node.js and OpenTelemetry Manual Instrumentation",description:"Quick start how to configure a Node.js app to use OpenTelemetry manual instrumentation with traces, and Tracetest for enhancing your E2E and integration tests with trace-based testing.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","nodejs","testing nodejs","nodejs observability","nodejs tracing","opentelemetry"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},i=void 0,o={id:"examples-tutorials/recipes/running-tracetest-without-a-trace-data-store-with-manual-instrumentation",title:"Node.js and OpenTelemetry Manual Instrumentation",description:"Quick start how to configure a Node.js app to use OpenTelemetry manual instrumentation with traces, and Tracetest for enhancing your E2E and integration tests with trace-based testing.",source:"@site/docs/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store-with-manual-instrumentation.mdx",sourceDirName:"examples-tutorials/recipes",slug:"/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store-with-manual-instrumentation",permalink:"/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store-with-manual-instrumentation",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store-with-manual-instrumentation.mdx",tags:[],version:"current",frontMatter:{id:"running-tracetest-without-a-trace-data-store-with-manual-instrumentation",title:"Node.js and OpenTelemetry Manual Instrumentation",description:"Quick start how to configure a Node.js app to use OpenTelemetry manual instrumentation with traces, and Tracetest for enhancing your E2E and integration tests with trace-based testing.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","nodejs","testing nodejs","nodejs observability","nodejs tracing","opentelemetry"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"guidesSidebar",previous:{title:"Node.js and OpenTelemetry",permalink:"/examples-tutorials/recipes/running-tracetest-without-a-trace-data-store"},next:{title:"Python with OpenTelemetry manual instrumention",permalink:"/examples-tutorials/recipes/running-python-app-with-opentelemetry-collector-and-tracetest"}},c={},l=[{value:"Node.js app with OpenTelemetry Manual Instrumentation and Tracetest",id:"nodejs-app-with-opentelemetry-manual-instrumentation-and-tracetest",level:2},{value:"Prerequisites",id:"prerequisites",level:2},{value:"Run This Quickstart Example",id:"run-this-quickstart-example",level:2},{value:"Project Structure",id:"project-structure",level:2},{value:"Configuring the Node.js App",id:"configuring-the-nodejs-app",level:2},{value:"Running the Node.js App and Tracetest",id:"running-the-nodejs-app-and-tracetest",level:2},{value:"Learn More",id:"learn-more",level:2}];function d(e){const t={a:"a",admonition:"admonition",code:"code",h2:"h2",img:"img",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,a.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(t.admonition,{type:"note",children:(0,s.jsx)(t.p,{children:(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples/quick-start-nodejs-manual-instrumentation",children:"Check out the source code on GitHub here."})})}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://tracetest.io/",children:"Tracetest"})," is a testing tool based on ",(0,s.jsx)(t.a,{href:"https://opentelemetry.io/",children:"OpenTelemetry"})," that allows you to test your distributed application. It allows you to use data from distributed traces generated by OpenTelemetry to validate and assert if your application has the desired behavior defined by your test definitions."]}),"\n",(0,s.jsx)(t.h2,{id:"nodejs-app-with-opentelemetry-manual-instrumentation-and-tracetest",children:"Node.js app with OpenTelemetry Manual Instrumentation and Tracetest"}),"\n",(0,s.jsx)(t.p,{children:"This is a simple quick start on how to configure a Node.js app to use OpenTelemetry instrumentation with traces, and Tracetest for enhancing your E2E and integration tests with trace-based testing. This example includes manual instrumentation and a sample bookstore array that simulates fetching data from a database."}),"\n",(0,s.jsx)(t.h2,{id:"prerequisites",children:"Prerequisites"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"Tracetest Account"}),":"]}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsxs)(t.li,{children:["Sign up to ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io",children:(0,s.jsx)(t.code,{children:"app.tracetest.io"})})," or follow the ",(0,s.jsx)(t.a,{href:"/getting-started/overview",children:"get started"})," docs."]}),"\n",(0,s.jsxs)(t.li,{children:["Have access to the environment's ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io/retrieve-token",children:"agent API key"}),"."]}),"\n"]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"Docker"}),": Have ",(0,s.jsx)(t.a,{href:"https://docs.docker.com/get-docker/",children:"Docker"})," and ",(0,s.jsx)(t.a,{href:"https://docs.docker.com/compose/install/",children:"Docker Compose"})," installed on your machine."]}),"\n",(0,s.jsx)(t.h2,{id:"run-this-quickstart-example",children:"Run This Quickstart Example"}),"\n",(0,s.jsx)(t.p,{children:"The example below is provided as part of the Tracetest project. You can download and run the example by following these steps:"}),"\n",(0,s.jsx)(t.p,{children:"Clone the Tracetest project and go to the Node.js Quickstart with Manual Instrumentation:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"git clone https://github.com/kubeshop/tracetest\ncd tracetest/examples/quick-start-nodejs-manual-instrumentation\n"})}),"\n",(0,s.jsx)(t.p,{children:"Follow these instructions to run the quick start:"}),"\n",(0,s.jsxs)(t.ol,{children:["\n",(0,s.jsxs)(t.li,{children:["Copy the ",(0,s.jsx)(t.code,{children:".env.template"})," file to ",(0,s.jsx)(t.code,{children:".env"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Fill out the ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io/retrieve-token",children:"TRACETEST_TOKEN and ENVIRONMENT_ID"})," details by editing your ",(0,s.jsx)(t.code,{children:".env"})," file."]}),"\n",(0,s.jsxs)(t.li,{children:["Run ",(0,s.jsx)(t.code,{children:"docker compose run tracetest-run"}),"."]}),"\n",(0,s.jsx)(t.li,{children:"Follow the links in the output to view the test results."}),"\n"]}),"\n",(0,s.jsx)(t.p,{children:"Follow along with the sections below for a detailed breakdown of what the example you just ran did and how it works."}),"\n",(0,s.jsx)(t.h2,{id:"project-structure",children:"Project Structure"}),"\n",(0,s.jsx)(t.p,{children:"The quick start Node.js project is built with Docker Compose and contains the Tracetest Agent and a Node.js app."}),"\n",(0,s.jsxs)(t.p,{children:["The ",(0,s.jsx)(t.code,{children:"docker-compose.yaml"})," file in the root directory of the quick start runs the Node.js app and the ",(0,s.jsx)(t.a,{href:"/concepts/agent",children:"Tracetest Agent"})," setup."]}),"\n",(0,s.jsx)(t.h2,{id:"configuring-the-nodejs-app",children:"Configuring the Node.js App"}),"\n",(0,s.jsxs)(t.p,{children:["The Node.js app is a simple Express app, contained in the ",(0,s.jsx)(t.code,{children:"app.js"})," file."]}),"\n",(0,s.jsxs)(t.p,{children:["Configure the ",(0,s.jsx)(t.code,{children:".env"})," like shown below."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:'# Get the required information here: https://app.tracetest.io/retrieve-token\n\nTRACETEST_TOKEN="<YOUR_TRACETEST_TOKEN>"\nTRACETEST_ENVIRONMENT_ID="<YOUR_ENV_ID>"\n\n# GRPC\nOTEL_EXPORTER_OTLP_TRACES_ENDPOINT="http://tracetest-agent:4317/"\n# or, use HTTP\n# OTEL_EXPORTER_OTLP_TRACES_ENDPOINT="http://tracetest-agent:4318/v1/traces"\n'})}),"\n",(0,s.jsxs)(t.p,{children:["The OpenTelemetry tracing is contained in the ",(0,s.jsx)(t.code,{children:"tracing.otel.grpc.js"})," or ",(0,s.jsx)(t.code,{children:"tracing.otel.http.js"})," files. Traces will be sent to Tracetest Agent."]}),"\n",(0,s.jsxs)(t.p,{children:["Choosing the ",(0,s.jsx)(t.code,{children:"tracing.otel.grpc.js"})," file will send traces to Tracetest Agent's ",(0,s.jsx)(t.code,{children:"GRPC"})," endpoint."]}),"\n",(0,s.jsxs)(t.p,{children:["Enabling the tracer is done by preloading the trace file. As seen in the ",(0,s.jsx)(t.code,{children:"package.json"}),"."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-json",children:'"scripts": {\n  "app-with-grpc-tracer": "node -r ./tracing.otel.grpc.js app.js",\n  "availability-with-grpc-tracer": "node -r ./tracing.otel.grpc.js availability.js",\n},\n'})}),"\n",(0,s.jsx)(t.h2,{id:"running-the-nodejs-app-and-tracetest",children:"Running the Node.js App and Tracetest"}),"\n",(0,s.jsx)(t.p,{children:"To execute the tests, run this command:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"docker compose run tracetest-run\n"})}),"\n",(0,s.jsx)(t.p,{children:"This will:"}),"\n",(0,s.jsxs)(t.ol,{children:["\n",(0,s.jsx)(t.li,{children:"Start the Node.js app, the OpenTelemetry Collector, and send the traces to the Tracetest Agent."}),"\n",(0,s.jsx)(t.li,{children:"Start the Tracetest Agent."}),"\n",(0,s.jsx)(t.li,{children:"Configure the tracing backend and create tests in your environment."}),"\n",(0,s.jsx)(t.li,{children:"Run the tests."}),"\n"]}),"\n",(0,s.jsx)(t.p,{children:"The output of the test will look similar to this:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:'Configuring Tracetest\n SUCCESS  Successfully configured Tracetest CLI\nRunning Trace-Based Tests...\n\u2718 RunGroup: #gcMqU_jIR (https://app.tracetest.io/organizations/xxx/environments/xxx/run/gcMqU_jIR)\n Summary: 0 passed, 1 failed, 0 pending\n  \u2718 Books list with availability (https://app.tracetest.io/organizations/xxx/environments/xxx/test/phAZcrT4B/run/10/test) - trace id: 7b92f5f6633218bbaacf0b79fe9c8904\n\t\u2714 span[tracetest.span.type="http" name="GET /books" http.target="/books" http.method="GET"]\n\t\t\u2714 #b0dfd9d58a43a950\n\t\t\t\u2714 attr:tracetest.span.duration  < 500ms (30ms)\n\t\u2714 span[tracetest.span.type="general" name="Books List"]\n\t\t\u2714 #8155143d27cbb3ab\n\t\t\t\u2714 attr:books.list.count = 3 (3)\n\t\u2714 span[tracetest.span.type="http" name="GET /availability/:bookId" http.method="GET"]\n\t\t\u2714 #1e2fcb5cefc171b5\n\t\t\t\u2714 attr:http.host = "availability:8080" (availability:8080)\n\t\t\u2714 #3e1fc86271bf1192\n\t\t\t\u2714 attr:http.host = "availability:8080" (availability:8080)\n\t\t\u2714 #9f29960c92ca268a\n\t\t\t\u2714 attr:http.host = "availability:8080" (availability:8080)\n\t\u2718 span[tracetest.span.type="general" name="Availablity check"]\n\t\t\u2714 #d191cefe48e65c74\n\t\t\t\u2714 attr:isAvailable = "true" (true)\n\t\t\u2714 #27989492723fd49a\n\t\t\t\u2714 attr:isAvailable = "true" (true)\n\t\t\u2718 #89e0d3b186bc7f0d\n\t\t\t\u2718 attr:isAvailable = "true" (false) (https://app.tracetest.io/organizations/xxx/environments/xxx/test/phAZcrT4B/run/10/test?selectedAssertion=3&selectedSpan=89e0d3b186bc7f0d)\n\n\t\u2718 Required gates\n\t\t\u2718 test-specs\n'})}),"\n",(0,s.jsx)(t.p,{children:(0,s.jsx)(t.img,{src:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1715607148/docs/app.tracetest.io_organizations_ttorg_e66318ba6544b856_environments_ttenv_56b66bc6e1a1cbbd_test_e9CCQuLSg_run_3_selectedAssertion_0_selectedSpan_d0c03aa5d02b9975_uqhhwl.png",alt:"assertion"})}),"\n",(0,s.jsx)(t.h2,{id:"learn-more",children:"Learn More"}),"\n",(0,s.jsxs)(t.p,{children:["Feel free to check out our ",(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples",children:"examples in GitHub"})," and join our ",(0,s.jsx)(t.a,{href:"https://dub.sh/tracetest-community",children:"Slack Community"})," for more info!"]})]})}function h(e={}){const{wrapper:t}={...(0,a.R)(),...e.components};return t?(0,s.jsx)(t,{...e,children:(0,s.jsx)(d,{...e})}):d(e)}},28453:(e,t,n)=>{n.d(t,{R:()=>i,x:()=>o});var s=n(96540);const a={},r=s.createContext(a);function i(e){const t=s.useContext(r);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function o(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(a):e.components||a:i(e.components),s.createElement(r.Provider,{value:t},e.children)}}}]);