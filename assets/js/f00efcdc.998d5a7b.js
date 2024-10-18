"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[4554],{88328:(e,t,s)=>{s.r(t),s.d(t,{assets:()=>c,contentTitle:()=>a,default:()=>d,frontMatter:()=>n,metadata:()=>i,toc:()=>h});var r=s(74848),o=s(28453);const n={id:"running-tracetest-with-aws-x-ray-pokeshop",title:"Pokeshop API with AWS X-Ray (Node.js SDK) and AWS Distro for OpenTelemetry",description:"Quick start on how to configure the Pokeshop API Demo with OpenTelemetry traces, AWS X-Ray as a trace data store, and Tracetest for enhancing your E2E and integration tests with trace-based testing.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","aws","aws observability","aws tracing","aws xray","aws xray tracing","aws adot collector","adot collector","opentelemetry"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},a=void 0,i={id:"examples-tutorials/recipes/running-tracetest-with-aws-x-ray-pokeshop",title:"Pokeshop API with AWS X-Ray (Node.js SDK) and AWS Distro for OpenTelemetry",description:"Quick start on how to configure the Pokeshop API Demo with OpenTelemetry traces, AWS X-Ray as a trace data store, and Tracetest for enhancing your E2E and integration tests with trace-based testing.",source:"@site/docs/examples-tutorials/recipes/running-tracetest-with-aws-x-ray-pokeshop.mdx",sourceDirName:"examples-tutorials/recipes",slug:"/examples-tutorials/recipes/running-tracetest-with-aws-x-ray-pokeshop",permalink:"/examples-tutorials/recipes/running-tracetest-with-aws-x-ray-pokeshop",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/examples-tutorials/recipes/running-tracetest-with-aws-x-ray-pokeshop.mdx",tags:[],version:"current",frontMatter:{id:"running-tracetest-with-aws-x-ray-pokeshop",title:"Pokeshop API with AWS X-Ray (Node.js SDK) and AWS Distro for OpenTelemetry",description:"Quick start on how to configure the Pokeshop API Demo with OpenTelemetry traces, AWS X-Ray as a trace data store, and Tracetest for enhancing your E2E and integration tests with trace-based testing.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","aws","aws observability","aws tracing","aws xray","aws xray tracing","aws adot collector","adot collector","opentelemetry"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"guidesSidebar",previous:{title:"Node.js with AWS X-Ray (Node.js SDK) and AWS Distro for OpenTelemetry",permalink:"/examples-tutorials/recipes/running-tracetest-with-aws-x-ray-adot"},next:{title:".NET Step Functions with AWS X-Ray, AWS Distro for OpenTelemetry, and Terraform",permalink:"/examples-tutorials/recipes/running-tracetest-with-step-functions-terraform"}},c={},h=[{value:"Pokeshop API Demo with AWS X-Ray (Node.js SDK), AWS Distro for OpenTelemetry and Tracetest",id:"pokeshop-api-demo-with-aws-x-ray-nodejs-sdk-aws-distro-for-opentelemetry-and-tracetest",level:2},{value:"Prerequisites",id:"prerequisites",level:2},{value:"Run This Quckstart Example",id:"run-this-quckstart-example",level:2},{value:"Project Structure",id:"project-structure",level:2},{value:"Configuring the Pokeshop Demo App",id:"configuring-the-pokeshop-demo-app",level:2},{value:"Run the Pokeshop Demo App, ADOT, and Tracetest Agent with Docker Compose",id:"run-the-pokeshop-demo-app-adot-and-tracetest-agent-with-docker-compose",level:2},{value:"Learn More",id:"learn-more",level:2}];function l(e){const t={a:"a",admonition:"admonition",code:"code",h2:"h2",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,o.R)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(t.admonition,{type:"note",children:(0,r.jsx)(t.p,{children:(0,r.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-amazon-x-ray-pokeshop",children:"Check out the source code on GitHub here."})})}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.a,{href:"https://tracetest.io/",children:"Tracetest"})," is a testing tool based on ",(0,r.jsx)(t.a,{href:"https://opentelemetry.io/",children:"OpenTelemetry"})," that allows you to test your distributed application. It allows you to use data from distributed traces generated by OpenTelemetry to validate and assert if your application has the desired behavior defined by your test definitions."]}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.a,{href:"https://aws.amazon.com/xray/",children:"AWS X-Ray"})," provides a complete view of requests as they travel through your application and filters visual data across payloads, functions, traces, services, APIs, and more with no-code and low-code motions."]}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.a,{href:"https://aws-otel.github.io/docs/getting-started/collector",children:"AWS Distro for OpenTelemetry (ADOT)"})," is a secure, production-ready, AWS-supported distribution of the OpenTelemetry project. Part of the Cloud Native Computing Foundation, OpenTelemetry provides open source APIs, libraries, and agents to collect distributed traces and metrics for application monitoring."]}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.a,{href:"https://docs.tracetest.io/live-examples/pokeshop/overview",children:"Pokeshop API"})," As a testing ground, the team at Tracetest has implemented a sample instrumented API around the ",(0,r.jsx)(t.a,{href:"https://pokeapi.co/",children:"PokeAPI"}),"."]}),"\n",(0,r.jsx)(t.h2,{id:"pokeshop-api-demo-with-aws-x-ray-nodejs-sdk-aws-distro-for-opentelemetry-and-tracetest",children:"Pokeshop API Demo with AWS X-Ray (Node.js SDK), AWS Distro for OpenTelemetry and Tracetest"}),"\n",(0,r.jsx)(t.p,{children:"This is a simple quick start guide on how to configure a fully instrumented API to be used with Tracetest for enhancing your E2E and integration tests with trace-based testing. The infrastructure will use AWS X-Ray as the trace data store, the ADOT as a middleware and the Pokeshop API to generate the telemetry data."}),"\n",(0,r.jsx)(t.h2,{id:"prerequisites",children:"Prerequisites"}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.strong,{children:"Tracetest Account"}),":"]}),"\n",(0,r.jsxs)(t.ul,{children:["\n",(0,r.jsxs)(t.li,{children:["Sign up to ",(0,r.jsx)(t.a,{href:"https://app.tracetest.io",children:(0,r.jsx)(t.code,{children:"app.tracetest.io"})})," or follow the ",(0,r.jsx)(t.a,{href:"/getting-started/overview",children:"get started"})," docs."]}),"\n",(0,r.jsxs)(t.li,{children:["Have access to the environment's ",(0,r.jsx)(t.a,{href:"https://app.tracetest.io/retrieve-token",children:"agent API key"}),"."]}),"\n"]}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.strong,{children:"AWS Account"}),":"]}),"\n",(0,r.jsxs)(t.ul,{children:["\n",(0,r.jsxs)(t.li,{children:["Sign up to ",(0,r.jsx)(t.a,{href:"https://aws.amazon.com/",children:"AWS"}),"."]}),"\n",(0,r.jsxs)(t.li,{children:["Install ",(0,r.jsx)(t.a,{href:"https://docs.aws.amazon.com/cli/latest/userguide/getting-started-install.html",children:"AWS CLI"}),"."]}),"\n"]}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.strong,{children:"Docker"}),": Have ",(0,r.jsx)(t.a,{href:"https://docs.docker.com/get-docker/",children:"Docker"})," and ",(0,r.jsx)(t.a,{href:"https://docs.docker.com/compose/install/",children:"Docker Compose"})," installed on your machine."]}),"\n",(0,r.jsx)(t.h2,{id:"run-this-quckstart-example",children:"Run This Quckstart Example"}),"\n",(0,r.jsx)(t.p,{children:"The example below is provided as part of the Tracetest project. You can download and run the example by following these steps:"}),"\n",(0,r.jsx)(t.p,{children:"Clone the Tracetest project and go to the AWS X-Ray Pokeshop Quickstart:"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",children:"git clone https://github.com/kubeshop/tracetest\ncd tracetest/examples/tracetest-amazon-x-ray-pokeshop\n"})}),"\n",(0,r.jsx)(t.p,{children:"Follow these instructions to run the quick start:"}),"\n",(0,r.jsxs)(t.ol,{children:["\n",(0,r.jsxs)(t.li,{children:["Copy the ",(0,r.jsx)(t.code,{children:".env.template"})," file to ",(0,r.jsx)(t.code,{children:".env"}),"."]}),"\n",(0,r.jsxs)(t.li,{children:["Fill out the ",(0,r.jsx)(t.a,{href:"https://app.tracetest.io/retrieve-token",children:"TRACETEST_TOKEN and ENVIRONMENT_ID"})," details by editing your ",(0,r.jsx)(t.code,{children:".env"})," file."]}),"\n",(0,r.jsxs)(t.li,{children:["Fill out the AWS credentials in the ",(0,r.jsx)(t.code,{children:".env"})," file. You can ",(0,r.jsx)(t.a,{href:"/configuration/connecting-to-data-stores/awsxray",children:"create credentials by following this guide"}),"."]}),"\n",(0,r.jsxs)(t.li,{children:["Run ",(0,r.jsx)(t.code,{children:"docker compose run tracetest-run"}),"."]}),"\n",(0,r.jsx)(t.li,{children:"Follow the links in the output to view the test results."}),"\n"]}),"\n",(0,r.jsx)(t.p,{children:"Follow along with the sections below for an in detail breakdown of what the example you just ran did and how it works."}),"\n",(0,r.jsx)(t.h2,{id:"project-structure",children:"Project Structure"}),"\n",(0,r.jsxs)(t.p,{children:["The project contains the ",(0,r.jsx)(t.a,{href:"/getting-started/install-agent",children:"Tracetest Agent"}),", and the ",(0,r.jsx)(t.a,{href:"/live-examples/pokeshop/overview",children:"Pokeshop Demo app"}),"."]}),"\n",(0,r.jsxs)(t.p,{children:["The ",(0,r.jsx)(t.code,{children:"docker-compose.yaml"})," file in the root directory of the quick start runs the Pokeshop Demo app, ",(0,r.jsx)(t.a,{href:"https://aws-otel.github.io/docs/getting-started/collector",children:"AWS Distro for OpenTelemetry (ADOT)"}),", and the ",(0,r.jsx)(t.a,{href:"/concepts/agent",children:"Tracetest Agent"})," setup."]}),"\n",(0,r.jsx)(t.h2,{id:"configuring-the-pokeshop-demo-app",children:"Configuring the Pokeshop Demo App"}),"\n",(0,r.jsx)(t.p,{children:"The Pokeshop API is a fully instrumented REST API that makes use of different services to mimic a real life scenario."}),"\n",(0,r.jsxs)(t.p,{children:["It is instrumented using the ",(0,r.jsx)(t.a,{href:"https://opentelemetry.io/docs/instrumentation/js/getting-started/nodejs/",children:"OpenTelemetry standards for Node.js"}),", sending the data to the ADOT collector that will be pushing the telemetry information to both the AWS X-Ray service."]}),"\n",(0,r.jsxs)(t.p,{children:["This is a ",(0,r.jsx)(t.a,{href:"https://github.com/kubeshop/pokeshop/blob/master/api/src/telemetry/tracing.ts",children:"fragment from the main tracing file"})," from the ",(0,r.jsx)(t.a,{href:"https://github.com/kubeshop/pokeshop",children:"Pokeshop API repo"}),"."]}),"\n",(0,r.jsxs)(t.p,{children:["Configure the ",(0,r.jsx)(t.code,{children:".env"})," like shown below."]}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",children:'# Get the required information here: https://app.tracetest.io/retrieve-token\n\nTRACETEST_TOKEN="<YOUR_TRACETEST_TOKEN>"\nTRACETEST_ENVIRONMENT_ID="<YOUR_ENV_ID>"\n\nAWS_ACCESS_KEY_ID="<YOUR_AWS_ACCESS_KEY_ID>"\nAWS_SECRET_ACCESS_KEY="<YOUR_AWS_SECRET_ACCESS_KEY>"\nAWS_SESSION_TOKEN="<YOURAWS_SESSION_TOKEN>"\nAWS_REGION="<YOUR_AWS_REGION>"\n'})}),"\n",(0,r.jsx)(t.h2,{id:"run-the-pokeshop-demo-app-adot-and-tracetest-agent-with-docker-compose",children:"Run the Pokeshop Demo App, ADOT, and Tracetest Agent with Docker Compose"}),"\n",(0,r.jsxs)(t.p,{children:["The ",(0,r.jsxs)(t.a,{href:"https://github.com/kubeshop/tracetest/blob/main/examples/tracetest-amazon-x-ray-pokeshop/docker-compose.yaml",children:[(0,r.jsx)(t.code,{children:"docker-compose.yaml"})," file"]})," in the root directory contains the Pokeshop Demo app services."]}),"\n",(0,r.jsxs)(t.p,{children:["The ",(0,r.jsxs)(t.a,{href:"https://github.com/kubeshop/tracetest/blob/main/examples/tracetest-amazon-x-ray-pokeshop/docker-compose.yaml",children:[(0,r.jsx)(t.code,{children:"docker-compose.yaml"})," file"]})," also contains the Tracetest Agent and ADOT."]}),"\n",(0,r.jsx)(t.p,{children:"To run everything including Tracetest tests, run this command:"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",children:"docker compose run tracetest-run\n"})}),"\n",(0,r.jsx)(t.p,{children:"This will:"}),"\n",(0,r.jsxs)(t.ol,{children:["\n",(0,r.jsx)(t.li,{children:"Start the Pokeshop app, the OpenTelemetry Collector, and send the traces to AWS X-Ray."}),"\n",(0,r.jsx)(t.li,{children:"Start the Tracetest Agent."}),"\n",(0,r.jsx)(t.li,{children:"Configure the tracing backend and create tests in your environment."}),"\n",(0,r.jsx)(t.li,{children:"Run the tests."}),"\n"]}),"\n",(0,r.jsx)(t.h2,{id:"learn-more",children:"Learn More"}),"\n",(0,r.jsxs)(t.p,{children:["Please visit our ",(0,r.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples",children:"examples in GitHub"})," and join our ",(0,r.jsx)(t.a,{href:"https://dub.sh/tracetest-community",children:"Slack Community"})," for more info!"]})]})}function d(e={}){const{wrapper:t}={...(0,o.R)(),...e.components};return t?(0,r.jsx)(t,{...e,children:(0,r.jsx)(l,{...e})}):l(e)}},28453:(e,t,s)=>{s.d(t,{R:()=>a,x:()=>i});var r=s(96540);const o={},n=r.createContext(o);function a(e){const t=r.useContext(n);return r.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function i(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:a(e.components),r.createElement(n.Provider,{value:t},e.children)}}}]);