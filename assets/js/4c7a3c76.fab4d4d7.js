"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[345],{72467:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>o,default:()=>h,frontMatter:()=>i,metadata:()=>a,toc:()=>l});var s=n(74848),r=n(28453);const i={id:"testing-vercel-functions-with-opentelemetry-tracetest",title:"Testing Vercel Functions (Next.js) with OpenTelemetry and Tracetest",description:"Quick start on how to configure Vercel functions with OpenTelemetry and Tracetest for enhancing your integration tests with trace-based testing.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","vercel","nextjs","end to end testing","end-to-end testing","integration testing","serverless testing","testing serverless functions","testing vercel functions","vercel testing","opentelemetry"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},o=void 0,a={id:"examples-tutorials/recipes/testing-vercel-functions-with-opentelemetry-tracetest",title:"Testing Vercel Functions (Next.js) with OpenTelemetry and Tracetest",description:"Quick start on how to configure Vercel functions with OpenTelemetry and Tracetest for enhancing your integration tests with trace-based testing.",source:"@site/docs/examples-tutorials/recipes/testing-vercel-functions-with-opentelemetry-tracetest.mdx",sourceDirName:"examples-tutorials/recipes",slug:"/examples-tutorials/recipes/testing-vercel-functions-with-opentelemetry-tracetest",permalink:"/examples-tutorials/recipes/testing-vercel-functions-with-opentelemetry-tracetest",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/examples-tutorials/recipes/testing-vercel-functions-with-opentelemetry-tracetest.mdx",tags:[],version:"current",frontMatter:{id:"testing-vercel-functions-with-opentelemetry-tracetest",title:"Testing Vercel Functions (Next.js) with OpenTelemetry and Tracetest",description:"Quick start on how to configure Vercel functions with OpenTelemetry and Tracetest for enhancing your integration tests with trace-based testing.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","vercel","nextjs","end to end testing","end-to-end testing","integration testing","serverless testing","testing serverless functions","testing vercel functions","vercel testing","opentelemetry"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"guidesSidebar",previous:{title:"Testing Kafka in a Go API with OpenTelemetry and Tracetest",permalink:"/examples-tutorials/recipes/testing-kafka-go-api-with-opentelemetry-tracetest"},next:{title:"Testing AWS Lambda Functions (Serverless Framework) with OpenTelemetry and Tracetest",permalink:"/examples-tutorials/recipes/testing-lambda-functions-with-opentelemetry-tracetest"}},c={},l=[{value:"Why is this important?",id:"why-is-this-important",level:2},{value:"Prerequisites",id:"prerequisites",level:2},{value:"Project Structure",id:"project-structure",level:2},{value:"1. Vercel (Next.js) Function",id:"1-vercel-nextjs-function",level:3},{value:"2. Tracetest",id:"2-tracetest",level:3},{value:"Docker Compose Network",id:"docker-compose-network",level:3},{value:"Vercel (Next.js) Function",id:"vercel-nextjs-function",level:2},{value:"Set up Environment Variables",id:"set-up-environment-variables",level:3},{value:"Start the Next.js Vercel Function",id:"start-the-nextjs-vercel-function",level:3},{value:"Testing the Vercel Function Locally",id:"testing-the-vercel-function-locally",level:2},{value:"Integration Testing the Vercel Function",id:"integration-testing-the-vercel-function",level:2},{value:"Learn More",id:"learn-more",level:2}];function d(e){const t={a:"a",admonition:"admonition",code:"code",h2:"h2",h3:"h3",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,r.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(t.admonition,{type:"note",children:(0,s.jsx)(t.p,{children:(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples/integration-testing-vercel-functions",children:"Check out the source code on GitHub here."})})}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://tracetest.io/",children:"Tracetest"})," is a testing tool based on ",(0,s.jsx)(t.a,{href:"https://opentelemetry.io/",children:"OpenTelemetry"})," that allows you to test your distributed application. It allows you to use data from distributed traces generated by OpenTelemetry to validate and assert if your application has the desired behavior defined by your test definitions."]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://vercel.com/",children:"Vercel"})," is a platform that hosts serverless functions and front-end code offering developers scalability and flexibility with no infrastructure overhead."]}),"\n",(0,s.jsx)(t.h2,{id:"why-is-this-important",children:"Why is this important?"}),"\n",(0,s.jsx)(t.p,{children:"Testing Serverless Functions has been a pain point for years. Not having visibility into the infrastructure and not knowing where a test fails causes the MTTR to be higher than for other tools. Including OpenTelemetry in Vercel functions exposes telemetry that you can use for both production visibility and trace-based testing."}),"\n",(0,s.jsxs)(t.p,{children:["This sample shows how to run integration tests against Vercel Functions using ",(0,s.jsx)(t.a,{href:"https://opentelemetry.io/",children:"OpenTelemetry"})," and Tracetest."]}),"\n",(0,s.jsx)(t.p,{children:"The Vercel function will fetch data from an external API, transform the data and insert it into a Vercel Postgres database. This particular flow has two failure points that are difficult to test."}),"\n",(0,s.jsxs)(t.ol,{children:["\n",(0,s.jsx)(t.li,{children:"Validating that an external API request from a Vercel function is successful."}),"\n",(0,s.jsx)(t.li,{children:"Validating that a Postgres insert request is successful."}),"\n"]}),"\n",(0,s.jsx)(t.h2,{id:"prerequisites",children:"Prerequisites"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"Tracetest Account"}),":"]}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsxs)(t.li,{children:["Sign up to ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io",children:(0,s.jsx)(t.code,{children:"app.tracetest.io"})})," or follow the ",(0,s.jsx)(t.a,{href:"/getting-started/overview",children:"get started"})," docs."]}),"\n",(0,s.jsxs)(t.li,{children:["Create an ",(0,s.jsx)(t.a,{href:"/concepts/environments",children:"environment"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Create an ",(0,s.jsx)(t.a,{href:"/concepts/environment-tokens",children:"environment token"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Have access to the environment's ",(0,s.jsx)(t.a,{href:"/configuration/agent",children:"agent API key"}),"."]}),"\n",(0,s.jsx)(t.li,{children:(0,s.jsx)(t.a,{href:"https://vercel.com/",children:"Vercel Account"})}),"\n",(0,s.jsx)(t.li,{children:(0,s.jsx)(t.a,{href:"https://vercel.com/docs/storage/vercel-postgres",children:"Vercel Postgres Database"})}),"\n"]}),"\n",(0,s.jsx)(t.p,{children:(0,s.jsx)(t.strong,{children:"Vercel Functions Example:"})}),"\n",(0,s.jsxs)(t.p,{children:["Clone the ",(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest",children:"Tracetest GitHub Repo"})," to your local machine, and open the Vercel example app."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"git clone https://github.com/kubeshop/tracetest.git\ncd tracetest/examples/integration-testing-vercel-functions\n"})}),"\n",(0,s.jsxs)(t.p,{children:["Before moving forward, run ",(0,s.jsx)(t.code,{children:"npm i"})," in the root folder to install the dependencies."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"npm i\n"})}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"Docker"}),":"]}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsxs)(t.li,{children:["Have ",(0,s.jsx)(t.a,{href:"https://docs.docker.com/get-docker/",children:"Docker"})," and ",(0,s.jsx)(t.a,{href:"https://docs.docker.com/compose/install/",children:"Docker Compose"})," installed on your machine."]}),"\n"]}),"\n",(0,s.jsx)(t.h2,{id:"project-structure",children:"Project Structure"}),"\n",(0,s.jsxs)(t.p,{children:["This is a ",(0,s.jsx)(t.a,{href:"https://nextjs.org/",children:"Next.js"})," project bootstrapped with ",(0,s.jsx)(t.a,{href:"https://github.com/vercel/next.js/tree/canary/packages/create-next-app",children:(0,s.jsx)(t.code,{children:"create-next-app"})}),"."]}),"\n",(0,s.jsxs)(t.p,{children:["It's using Vercel Functions via ",(0,s.jsx)(t.code,{children:"/pages/api"}),", with ",(0,s.jsx)(t.a,{href:"https://nextjs.org/docs/pages/building-your-application/optimizing/open-telemetry#manual-opentelemetry-configuration",children:"OpenTelemetry configured as explained in the Vercel docs"}),"."]}),"\n",(0,s.jsx)(t.h3,{id:"1-vercel-nextjs-function",children:"1. Vercel (Next.js) Function"}),"\n",(0,s.jsxs)(t.p,{children:["The ",(0,s.jsx)(t.code,{children:"docker-compose.yaml"})," file reference the Next.js app with ",(0,s.jsx)(t.code,{children:"next-app"}),"."]}),"\n",(0,s.jsx)(t.h3,{id:"2-tracetest",children:"2. Tracetest"}),"\n",(0,s.jsxs)(t.p,{children:["The ",(0,s.jsx)(t.code,{children:"docker-compose.yaml"})," file also has a Tracetest Agent service and an integration tests service."]}),"\n",(0,s.jsx)(t.h3,{id:"docker-compose-network",children:"Docker Compose Network"}),"\n",(0,s.jsxs)(t.p,{children:["All ",(0,s.jsx)(t.code,{children:"services"})," in the ",(0,s.jsx)(t.code,{children:"docker-compose.yaml"})," are on the same network and will be reachable by hostname from within other services. E.g. ",(0,s.jsx)(t.code,{children:"next-app:3000"})," in the ",(0,s.jsx)(t.code,{children:"test/api.pokemon.spec.docker.yaml"})," will map to the ",(0,s.jsx)(t.code,{children:"next-app"})," service."]}),"\n",(0,s.jsx)(t.h2,{id:"vercel-nextjs-function",children:"Vercel (Next.js) Function"}),"\n",(0,s.jsxs)(t.p,{children:["The Vercel Function is a simple API, ",(0,s.jsxs)(t.a,{href:"https://github.com/kubeshop/tracetest/blob/main/examples/integration-testing-vercel-functions/pages/api/pokemon.ts",children:["contained in the ",(0,s.jsx)(t.code,{children:"pages/api/pokemon.ts"})," file"]}),"."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-typescript",children:"import { trace, SpanStatusCode } from '@opentelemetry/api'\nimport type { NextApiRequest, NextApiResponse } from 'next'\nimport { sql } from '@vercel/postgres'\n\nexport async function addPokemon(pokemon: any) {\n  return await sql`\n    INSERT INTO pokemon (name)\n    VALUES (${pokemon.name})\n    RETURNING *;\n  `\n}\n\nexport async function getPokemon(pokemon: any) {\n  return await sql`\n    SELECT * FROM pokemon where id=${pokemon.id};\n  `\n}\n\nexport default async function handler(\n  req: NextApiRequest,\n  res: NextApiResponse\n) {\n  const activeSpan = trace.getActiveSpan()\n  const tracer = await trace.getTracer('integration-testing-vercel-functions')\n  \n  try {\n\n    const externalPokemon = await tracer.startActiveSpan('GET Pokemon from pokeapi.co', async (externalPokemonSpan) => {\n      const requestUrl = `https://pokeapi.co/api/v2/pokemon/${req.body.id || '6'}`\n      const response = await fetch(requestUrl)\n      const { id, name } = await response.json()\n\n      externalPokemonSpan.setStatus({ code: SpanStatusCode.OK, message: String(\"Pokemon fetched successfully!\") })\n      externalPokemonSpan.setAttribute('pokemon.name', name)\n      externalPokemonSpan.setAttribute('pokemon.id', id)\n      externalPokemonSpan.end()\n\n      return { id, name }\n    })\n\n    const addedPokemon = await tracer.startActiveSpan('Add Pokemon to Vercel Postgres', async (addedPokemonSpan) => {\n      const { rowCount, rows: [addedPokemon, ...rest] } = await addPokemon(externalPokemon)\n      addedPokemonSpan.setAttribute('pokemon.isAdded', rowCount === 1)\n      addedPokemonSpan.setAttribute('pokemon.added.name', addedPokemon.name)\n      addedPokemonSpan.end()\n      return addedPokemon\n    })\n    \n    res.status(200).json(addedPokemon)\n\n  } catch (err) {\n    activeSpan?.setAttribute('error', String(err))\n    activeSpan?.recordException(String(err))\n    activeSpan?.setStatus({ code: SpanStatusCode.ERROR, message: String(err) })\n    res.status(500).json({ error: 'failed to load data' })\n  } finally {\n    activeSpan?.end()\n  }\n}\n"})}),"\n",(0,s.jsxs)(t.p,{children:["The OpenTelemetry tracing is ",(0,s.jsxs)(t.a,{href:"https://github.com/kubeshop/tracetest/blob/main/examples/integration-testing-vercel-functions/instrumentation.node.ts",children:["contained in the ",(0,s.jsx)(t.code,{children:"instrumentation.node.ts"})," file"]}),". Traces will be sent to the Tracetest Agent."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-typescript",children:"import { NodeSDK } from '@opentelemetry/sdk-node'\nimport { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http'\nimport { Resource } from '@opentelemetry/resources'\nimport { SemanticResourceAttributes } from '@opentelemetry/semantic-conventions'\nimport { getNodeAutoInstrumentations } from '@opentelemetry/auto-instrumentations-node'\nimport { FetchInstrumentation } from '@opentelemetry/instrumentation-fetch'\n\nconst sdk = new NodeSDK({\n  // The OTEL_EXPORTER_OTLP_ENDPOINT env var is passed into \"new OTLPTraceExporter\" automatically.\n  // If the OTEL_EXPORTER_OTLP_ENDPOINT env var is not set the \"new OTLPTraceExporter\" will\n  // default to use \"http://localhost:4317\" for gRPC and \"http://localhost:4318\" for HTTP.\n  // This sample is using HTTP.\n  traceExporter: new OTLPTraceExporter(),\n  instrumentations: [\n    getNodeAutoInstrumentations(),\n    new FetchInstrumentation(),\n  ],\n  resource: new Resource({\n    [SemanticResourceAttributes.SERVICE_NAME]: 'integration-testing-vercel-functions',\n  }),\n})\nsdk.start()\n"})}),"\n",(0,s.jsx)(t.h3,{id:"set-up-environment-variables",children:"Set up Environment Variables"}),"\n",(0,s.jsxs)(t.p,{children:["Edit the ",(0,s.jsx)(t.code,{children:".env.development"})," file. Add your Vercel Postgres env vars."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",metastring:"title=.env.development.local",children:'# OTLP HTTP\nOTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:4318"\n\n# Vercel Postgres\nPOSTGRES_DATABASE="**********"\nPOSTGRES_HOST="**********"\nPOSTGRES_PASSWORD="**********"\nPOSTGRES_PRISMA_URL="**********"\nPOSTGRES_URL="**********"\nPOSTGRES_URL_NON_POOLING="**********"\nPOSTGRES_USER="**********"\n'})}),"\n",(0,s.jsx)(t.h3,{id:"start-the-nextjs-vercel-function",children:"Start the Next.js Vercel Function"}),"\n",(0,s.jsx)(t.p,{children:"Spin up your Next.js app."}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"npm run dev\n"})}),"\n",(0,s.jsxs)(t.p,{children:["This starts the function on ",(0,s.jsx)(t.code,{children:"http://localhost:3000/api/pokemon"}),"."]}),"\n",(0,s.jsx)(t.h2,{id:"testing-the-vercel-function-locally",children:"Testing the Vercel Function Locally"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"/getting-started/install-cli",children:"Download the CLI"})," for your operating system."]}),"\n",(0,s.jsxs)(t.p,{children:["The CLI is bundled with ",(0,s.jsx)(t.a,{href:"/concepts/agent/",children:"Tracetest Agent"})," that runs in your infrastructure to collect responses and traces for tests."]}),"\n",(0,s.jsxs)(t.p,{children:["To start Tracetest Agent add the ",(0,s.jsx)(t.code,{children:"--api-key"})," from your environment."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",metastring:"title=Terminal",children:"tracetest start --api-key YOUR_AGENT_API_KEY\n"})}),"\n",(0,s.jsxs)(t.p,{children:["Run a test with the test definition ",(0,s.jsx)(t.code,{children:"test/api.pokemon.spec.development.yaml"}),"."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-yaml",metastring:"title=test/api.pokemon.spec.development.yaml",children:'type: Test\nspec:\n  id: kv8C-hOSR\n  name: Test API\n  trigger:\n    type: http\n    httpRequest:\n      method: POST\n      url: http://localhost:3000/api/pokemon\n      body: "{\\n  \\"id\\": \\"6\\"\\n}"\n      headers:\n      - key: Content-Type\n        value: application/json\n  specs:\n  - selector: span[tracetest.span.type="http"]\n    name: "All HTTP Spans: Status  code is 200"\n    assertions:\n    - attr:http.status_code = 200\n'})}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",metastring:"title=Terminal",children:"tracetest run test -f ./test/api.pokemon.spec.development.yaml --required-gates test-specs --output pretty\n\n[Output]\n\u2714 Test API (https://app.tracetest.io/organizations/<YOUR_ORG>/environments/<YOUR_ENV>/test/-gjd4idIR/run/22/test) - trace id: f2250362ff2f70f8f5be7b2fba74e4b2\n    \u2714 All HTTP Spans: Status code is 200\n"})}),"\n",(0,s.jsx)(t.h2,{id:"integration-testing-the-vercel-function",children:"Integration Testing the Vercel Function"}),"\n",(0,s.jsxs)(t.p,{children:["Edit the ",(0,s.jsx)(t.code,{children:".env.docker"})," file to use your Vercel Postgres env vars."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",metastring:"title=.env.docker",children:'# OTLP HTTP\nOTEL_EXPORTER_OTLP_ENDPOINT="http://tracetest-agent:4318"\n\n# Vercel Postgres\nPOSTGRES_DATABASE="**********"\nPOSTGRES_HOST="**********"\nPOSTGRES_PASSWORD="**********"\nPOSTGRES_PRISMA_URL="**********"\nPOSTGRES_URL="**********"\nPOSTGRES_URL_NON_POOLING="**********"\nPOSTGRES_USER="**********"\n'})}),"\n",(0,s.jsxs)(t.p,{children:["This configures the ",(0,s.jsx)(t.code,{children:"OTEL_EXPORTER_OTLP_ENDPOINT"})," to send traces to Tracetest Agent."]}),"\n",(0,s.jsxs)(t.p,{children:["Edit the ",(0,s.jsx)(t.code,{children:"docker-compose.yaml"})," in the root directory. Add your ",(0,s.jsx)(t.code,{children:"TRACETEST_API_KEY"}),"."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-yaml",metastring:"title=docker-compose.yaml",children:"  # [...]\n  tracetest-agent:\n    image: kubeshop/tracetest-agent:latest\n    environment:\n      - TRACETEST_API_KEY=YOUR_TRACETEST_API_KEY # Find the Agent API Key here: https://docs.tracetest.io/configuration/agent\n    ports:\n      - 4317:4317\n      - 4318:4318\n    networks:\n      - tracetest\n"})}),"\n",(0,s.jsxs)(t.p,{children:["Edit the ",(0,s.jsx)(t.code,{children:"run.bash"}),". Add your ",(0,s.jsx)(t.code,{children:"TRACETEST_API_TOKEN"}),"."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"#/bin/bash\n\n# Find the API Token here: https://docs.tracetest.io/concepts/environment-tokens\ntracetest configure -t YOUR_TRACETEST_API_TOKEN\ntracetest run test -f ./api.pokemon.spec.docker.yaml --required-gates test-specs --output pretty\n"})}),"\n",(0,s.jsx)(t.p,{children:"Now you can run the Vercel function and Tracetest Agent!"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"docker compose up -d --build\n"})}),"\n",(0,s.jsx)(t.p,{children:"And, trigger the integration tests."}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"docker compose run integration-tests\n\n[Ouput]\n[+] Creating 1/0\n \u2714 Container integration-testing-vercel-functions-tracetest-agent-1  Running                                                                                             0.0s\n SUCCESS  Successfully configured Tracetest CLI\n\u2714 Test API (https://app.tracetest.io/organizations/<YOUR_ORG>/environments/<YOUR_ENV>/test/p00W82OIR/run/8/test) - trace id: d64ab3a6f52a98141d26679fff3373b6\n    \u2714 All HTTP Spans: Status code is 200\n"})}),"\n",(0,s.jsx)(t.h2,{id:"learn-more",children:"Learn More"}),"\n",(0,s.jsxs)(t.p,{children:["Feel free to check out our ",(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples",children:"examples in GitHub"})," and join our ",(0,s.jsx)(t.a,{href:"https://dub.sh/tracetest-community",children:"Slack Community"})," for more info!"]})]})}function h(e={}){const{wrapper:t}={...(0,r.R)(),...e.components};return t?(0,s.jsx)(t,{...e,children:(0,s.jsx)(d,{...e})}):d(e)}},28453:(e,t,n)=>{n.d(t,{R:()=>o,x:()=>a});var s=n(96540);const r={},i=s.createContext(r);function o(e){const t=s.useContext(i);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function a(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:o(e.components),s.createElement(i.Provider,{value:t},e.children)}}}]);