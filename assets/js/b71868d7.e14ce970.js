"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[2376],{26019:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>o,default:()=>d,frontMatter:()=>r,metadata:()=>a,toc:()=>l});var s=n(74848),i=n(28453);const r={id:"typescript",title:"Programmatically triggered trace-based tests using Tracetest and TypeScript",description:"Quickstart on how to use the Tracetest NPM @tracetest/client Typescript package to Programatically trigger Trace-Based Tests.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","typescript","programmatically"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},o=void 0,a={id:"tools-and-integrations/typescript",title:"Programmatically triggered trace-based tests using Tracetest and TypeScript",description:"Quickstart on how to use the Tracetest NPM @tracetest/client Typescript package to Programatically trigger Trace-Based Tests.",source:"@site/docs/tools-and-integrations/typescript.mdx",sourceDirName:"tools-and-integrations",slug:"/tools-and-integrations/typescript",permalink:"/tools-and-integrations/typescript",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/tools-and-integrations/typescript.mdx",tags:[],version:"current",frontMatter:{id:"typescript",title:"Programmatically triggered trace-based tests using Tracetest and TypeScript",description:"Quickstart on how to use the Tracetest NPM @tracetest/client Typescript package to Programatically trigger Trace-Based Tests.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing","typescript","programmatically"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"guidesSidebar",previous:{title:"Trace-Based End to End Testing with Playwright and Tracetest",permalink:"/tools-and-integrations/playwright"},next:{title:"Running Tracetest with Keptn",permalink:"/tools-and-integrations/keptn"}},c={},l=[{value:"Why is this important?",id:"why-is-this-important",level:2},{value:"The <code>@tracetest/client</code> NPM Package",id:"the-tracetestclient-npm-package",level:2},{value:"How It Works",id:"how-it-works",level:2},{value:"Today You&#39;ll Learn How to integrate Trace-Based Tests with your Typescript Code",id:"today-youll-learn-how-to-integrate-trace-based-tests-with-your-typescript-code",level:2},{value:"Requirements",id:"requirements",level:2},{value:"Run This Quckstart Example",id:"run-this-quckstart-example",level:2},{value:"Project Structure",id:"project-structure",level:2},{value:"Installing the <code>@tracetest/client</code> NPM Package",id:"installing-the-tracetestclient-npm-package",level:2},{value:"Tracetest Test Definitions",id:"tracetest-test-definitions",level:2},{value:"Creating the Typescript Script",id:"creating-the-typescript-script",level:2},{value:"Running the Full Example",id:"running-the-full-example",level:2},{value:"Finding the Results",id:"finding-the-results",level:2},{value:"What&#39;s Next?",id:"whats-next",level:2},{value:"Learn More",id:"learn-more",level:2}];function p(e){const t={a:"a",admonition:"admonition",code:"code",h2:"h2",img:"img",li:"li",mermaid:"mermaid",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,i.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(t.admonition,{title:"Commercial Feature",type:"tip",children:(0,s.jsx)(t.p,{children:(0,s.jsx)(t.a,{href:"https://tracetest.io/pricing",children:"Feature available only in Cloud-based Managed Tracetest & Enterprise Self-hosted Tracetest."})})}),"\n",(0,s.jsx)(t.admonition,{type:"note",children:(0,s.jsx)(t.p,{children:(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples/quick-start-typescript",children:"Check out the source code on GitHub here."})})}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.a,{href:"https://app.tracetest.io/",children:"Tracetest"})," is a testing tool based on ",(0,s.jsx)(t.a,{href:"https://opentelemetry.io/",children:"OpenTelemetry"})," that allows you to test your distributed application. It allows you to use data from distributed traces generated by OpenTelemetry to validate and assert if your application has the desired behavior defined by your test definitions."]}),"\n",(0,s.jsx)(t.p,{children:"JavaScript/TypeScript is today the most popular language for web development, and it is also the most popular language for writing tests and automation scripts."}),"\n",(0,s.jsx)(t.h2,{id:"why-is-this-important",children:"Why is this important?"}),"\n",(0,s.jsxs)(t.p,{children:["When working with testing tools, the most important thing is to be able to integrate them into your existing workflow and tooling. This is why we have created the ",(0,s.jsx)(t.code,{children:"@tracetest/client"})," NPM package, which allows you to use the Tracetest platform to run trace-based tests from your existing JavaScript/TypeScript code.\nEnabling you to run tests at any point in your code, and not only at the end of the test run, allows you to use trace-based testing as a tool to help you develop your application."]}),"\n",(0,s.jsxs)(t.admonition,{type:"info",children:[(0,s.jsx)(t.p,{children:"Check out the hands-on tutorial on YouTube!"}),(0,s.jsx)("iframe",{width:"100%",height:"250",src:"https://www.youtube.com/embed/BOMjkiwyRzc",title:"YouTube video player",frameborder:"0",allow:"accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share",allowfullscreen:!0})]}),"\n",(0,s.jsxs)(t.h2,{id:"the-tracetestclient-npm-package",children:["The ",(0,s.jsx)(t.code,{children:"@tracetest/client"})," NPM Package"]}),"\n",(0,s.jsxs)(t.p,{children:["With the ",(0,s.jsxs)(t.a,{href:"https://www.npmjs.com/package/@tracetest/client",children:[(0,s.jsx)(t.code,{children:"@tracetest/client"})," NPM Package"]}),", you will unlock the power of OpenTelemetry that allows you to run deeper testing based on the traces and spans generated by each of the checkpoints that you define within your services."]}),"\n",(0,s.jsx)(t.h2,{id:"how-it-works",children:"How It Works"}),"\n",(0,s.jsxs)(t.p,{children:["The following is a high-level sequence diagram of how the Typescript script, the ",(0,s.jsx)(t.code,{children:"@tracetest/client"})," package, and Tracetest interact."]}),"\n",(0,s.jsx)(t.mermaid,{value:"sequenceDiagram\n  Typescript->>+Typescript: Execute tests\n  Typescript->>+@tracetest/client: Createst instance\n  @tracetest/client--\x3e>-Typescript: Ok\n  Typescript->>+Typescript: Run pre-steps\n  Typescript->>+@tracetest/client: Creates test\n  @tracetest/client--\x3e>-Typescript: Ok\n  Typescript->>+@tracetest/client: Runs test\n  @tracetest/client--\x3e>-Typescript: Ok\n  Typescript->>@tracetest/client: Waits for results and shows the summary"}),"\n",(0,s.jsx)(t.h2,{id:"today-youll-learn-how-to-integrate-trace-based-tests-with-your-typescript-code",children:"Today You'll Learn How to integrate Trace-Based Tests with your Typescript Code"}),"\n",(0,s.jsxs)(t.p,{children:["This is a simple quick-start guide on how to use the Tracetest ",(0,s.jsx)(t.code,{children:"@tracetest/client"})," NPM package to enhance your TypeScript toolkit with trace-based testing. The infrastructure will use the Pokeshop Demo as a testing ground, triggering requests against it and generating telemetry data."]}),"\n",(0,s.jsx)(t.h2,{id:"requirements",children:"Requirements"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"Tracetest Account"}),":"]}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsxs)(t.li,{children:["Sign up to ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io",children:(0,s.jsx)(t.code,{children:"app.tracetest.io"})})," or follow the ",(0,s.jsx)(t.a,{href:"/getting-started/overview",children:"get started"})," docs."]}),"\n",(0,s.jsxs)(t.li,{children:["Create an ",(0,s.jsx)(t.a,{href:"/concepts/environments",children:"environment"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Create an ",(0,s.jsx)(t.a,{href:"/concepts/environment-tokens",children:"environment token"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Have access to the environment's ",(0,s.jsx)(t.a,{href:"/configuration/agent",children:"agent API key"}),"."]}),"\n"]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"Docker"}),": Have ",(0,s.jsx)(t.a,{href:"https://docs.docker.com/get-docker/",children:"Docker"})," and ",(0,s.jsx)(t.a,{href:"https://docs.docker.com/compose/install/",children:"Docker Compose"})," installed on your machine."]}),"\n",(0,s.jsx)(t.h2,{id:"run-this-quckstart-example",children:"Run This Quckstart Example"}),"\n",(0,s.jsx)(t.p,{children:"The example below is provided as part of the Tracetest GitHub repository. You can download and run the example by following these steps:"}),"\n",(0,s.jsx)(t.p,{children:"Clone the Tracetest project and go to the TypeScript Quickstart:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"git clone https://github.com/kubeshop/tracetest\ncd tracetest/examples/quick-start-typescript\n"})}),"\n",(0,s.jsx)(t.p,{children:"Follow these instructions to run the included demo app and TypeScript example:"}),"\n",(0,s.jsxs)(t.ol,{children:["\n",(0,s.jsxs)(t.li,{children:["Copy the ",(0,s.jsx)(t.code,{children:".env.template"})," file to ",(0,s.jsx)(t.code,{children:".env"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Log into the ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io/",children:"Tracetest app"}),"."]}),"\n",(0,s.jsx)(t.li,{children:"This example is configured to use the OpenTelemetry Collector. Ensure the environment you will be using to run this example is also configured to use the OpenTelemetry Tracing Backend by clicking on Settings, Tracing Backend, OpenTelemetry, Save."}),"\n",(0,s.jsxs)(t.li,{children:["Fill out the ",(0,s.jsx)(t.a,{href:"https://docs.tracetest.io/concepts/environment-tokens",children:"token"})," and ",(0,s.jsx)(t.a,{href:"https://docs.tracetest.io/concepts/agent",children:"agent API key"})," details by editing your ",(0,s.jsx)(t.code,{children:".env"})," file. You can find these values in the Settings area for your environment."]}),"\n",(0,s.jsxs)(t.li,{children:["Run ",(0,s.jsx)(t.code,{children:"docker compose up -d"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Look for the ",(0,s.jsx)(t.code,{children:"tracetest-client"})," service in Docker and click on it to view the logs. It will show the results from the trace-based tests that are triggered from the ",(0,s.jsx)(t.code,{children:"index.ts"})," TypeScript file."]}),"\n",(0,s.jsx)(t.li,{children:"Follow the links in the log to to view the test runs programmatically created by your TypeScript test script."}),"\n"]}),"\n",(0,s.jsx)(t.p,{children:"Follow along with the sections below for an in detail breakdown of what the example you just ran did and how it works."}),"\n",(0,s.jsx)(t.h2,{id:"project-structure",children:"Project Structure"}),"\n",(0,s.jsx)(t.p,{children:"The quick start TypeScript project is built with Docker Compose."}),"\n",(0,s.jsxs)(t.p,{children:["The ",(0,s.jsx)(t.a,{href:"/live-examples/pokeshop/overview",children:"Pokeshop Demo App"})," is a complete example of a distributed application using different back-end and front-end services. We will be launching it and running tests against it as part of this example."]}),"\n",(0,s.jsxs)(t.p,{children:["The ",(0,s.jsx)(t.code,{children:"docker-compose.yaml"})," file in the root directory of the quick start runs the Pokeshop Demo app, the OpenTelemetry Collector setup, and the ",(0,s.jsx)(t.a,{href:"/concepts/agent",children:"Tracetest Agent"}),"."]}),"\n",(0,s.jsxs)(t.p,{children:["The TypeScript quick start has two primary files: a Typescript file ",(0,s.jsx)(t.code,{children:"definitions.ts"})," that defines two Tracetest tests, and a Typescript file 'index.ts' that imports these test definitions and uses the \"@tracetest/client\" NPM package to run them multiple times."]}),"\n",(0,s.jsx)(t.p,{children:"We will show you how to install the NPM package and use these two TypeScript files to programmatically run Tracetest tests."}),"\n",(0,s.jsxs)(t.h2,{id:"installing-the-tracetestclient-npm-package",children:["Installing the ",(0,s.jsx)(t.code,{children:"@tracetest/client"})," NPM Package"]}),"\n",(0,s.jsxs)(t.p,{children:["The first step when using the TypeScript NPM package is to install the ",(0,s.jsx)(t.code,{children:"@tracetest/client"})," NPM Package. It is as easy as running the following command:"]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"npm i @tracetest/client\n"})}),"\n",(0,s.jsxs)(t.p,{children:["Once you have installed the ",(0,s.jsx)(t.code,{children:"@tracetest/client"})," package, you can import it and start making use of it as any other library to trigger trace-based tests and run checks against the resulting telemetry data."]}),"\n",(0,s.jsx)(t.h2,{id:"tracetest-test-definitions",children:"Tracetest Test Definitions"}),"\n",(0,s.jsxs)(t.p,{children:["The ",(0,s.jsx)(t.code,{children:"definitions.ts"})," file contains the JSON version of the test definitions that will be used to run the tests. It uses the HTTP trigger to execute requests against the Pokeshop Demo."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-typescript",children:'import { TestResource } from "@tracetest/client/dist/modules/openapi-client";\n\nexport const importDefinition: TestResource = {\n  type: "Test",\n  spec: {\n    id: "99TOHzpSR",\n    name: "Typescript: Import a Pokemon",\n    trigger: {\n      type: "http",\n      httpRequest: {\n        method: "POST",\n        url: "${var:BASE_URL}/import",\n        body: \'{"id": ${var:POKEMON_ID}}\',\n        headers: [\n          {\n            key: "Content-Type",\n            value: "application/json",\n          },\n        ],\n      },\n    },\n    specs: [\n      {\n        selector: \'span[tracetest.span.type="general" name = "validate request"] span[tracetest.span.type="http"]\',\n        name: "All HTTP Spans: Status  code is 200",\n        assertions: ["attr:http.status_code = 200"],\n      },\n      {\n        selector: \'span[tracetest.span.type="http" name="GET" http.method="GET"]\',\n        assertions: [\'attr:http.route = "/api/v2/pokemon/${var:POKEMON_ID}"\'],\n      },\n      {\n        selector: \'span[tracetest.span.type="database"]\',\n        name: "All Database Spans: Processing time is less than 1s",\n        assertions: ["attr:tracetest.span.duration < 1s"],\n      },\n    ],\n    outputs: [\n      {\n        name: "DATABASE_POKEMON_ID",\n        selector:\n          \'span[tracetest.span.type="database" name="create pokeshop.pokemon" db.system="postgres" db.name="pokeshop" db.user="ashketchum" db.operation="create" db.sql.table="pokemon"]\',\n        value: "attr:db.result | json_path \'$.id\'",\n      },\n    ],\n  },\n};\n\nexport const deleteDefinition: TestResource = {\n  type: "Test",\n  spec: {\n    id: "C2gwdktIR",\n    name: "Typescript: Delete a Pokemon",\n    trigger: {\n      type: "http",\n      httpRequest: {\n        method: "DELETE",\n        url: "${var:BASE_URL}/${var:POKEMON_ID}",\n        headers: [\n          {\n            key: "Content-Type",\n            value: "application/json",\n          },\n        ],\n      },\n    },\n    specs: [\n      {\n        selector:\n          \'span[tracetest.span.type="database" db.system="redis" db.operation="del" db.redis.database_index="0"]\',\n        assertions: [\'attr:db.payload = \\\'{"key":"pokemon-${var:POKEMON_ID}"}\\\'\'],\n      },\n      {\n        selector:\n          \'span[tracetest.span.type="database" name="delete pokeshop.pokemon" db.system="postgres" db.name="pokeshop" db.user="ashketchum" db.operation="delete" db.sql.table="pokemon"]\',\n        assertions: ["attr:db.result = 1"],\n      },\n      {\n        selector: \'span[tracetest.span.type="database"]\',\n        name: "All Database Spans: Processing time is less than 100ms",\n        assertions: ["attr:tracetest.span.duration < 100ms"],\n      },\n    ],\n  },\n};\n'})}),"\n",(0,s.jsx)(t.h2,{id:"creating-the-typescript-script",children:"Creating the Typescript Script"}),"\n",(0,s.jsxs)(t.p,{children:["The ",(0,s.jsx)(t.code,{children:"index.ts"})," file contains the Typescript script that will be used to trigger requests against the Pokeshop Demo and run trace-based tests. The steps executed by this script are the following:"]}),"\n",(0,s.jsxs)(t.ol,{children:["\n",(0,s.jsxs)(t.li,{children:["Import the ",(0,s.jsx)(t.code,{children:"@tracetest/client"})," package."]}),"\n",(0,s.jsxs)(t.li,{children:["Create a new ",(0,s.jsx)(t.code,{children:"Tracetest"})," instance."]}),"\n",(0,s.jsxs)(t.li,{children:["Get the last imported Pokemon number from the ",(0,s.jsx)(t.code,{children:"GET /pokemon"})," endpoint using ",(0,s.jsx)(t.code,{children:"fetch"}),"."]}),"\n",(0,s.jsxs)(t.li,{children:["Import the following 5 Pokemon after the last number by triggering a trace-based test to the ",(0,s.jsx)(t.code,{children:"POST /import"})," endpoint."]}),"\n",(0,s.jsxs)(t.li,{children:["From each test output, get the ",(0,s.jsx)(t.code,{children:"DATABASE_POKEMON_ID"})," value and add it to a list."]}),"\n",(0,s.jsxs)(t.li,{children:["Delete the imported Pokemon by triggering a trace-based test to the ",(0,s.jsx)(t.code,{children:"DELETE /:id"})," endpoint."]}),"\n"]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-typescript",children:'import Tracetest from "@tracetest/client";\nimport { config } from "dotenv";\nimport { PokemonList } from "./types";\nimport { deleteDefinition, importDefinition } from "./definitions";\n\nconfig();\n\nconst { TRACETEST_API_TOKEN = "", POKESHOP_DEMO_URL = "http://api:8081" } = process.env;\n\nconst baseUrl = `${POKESHOP_DEMO_URL}/pokemon`;\n\nconst main = async () => {\n  const tracetest = await Tracetest(TRACETEST_API_TOKEN);\n\n  const getLastPokemonId = async (): Promise<number> => {\n    const response = await fetch(baseUrl);\n    const list = (await response.json()) as PokemonList;\n\n    return list.items.length + 1;\n  };\n\n  // get the initial pokemon from the API\n  const pokemonId = (await getLastPokemonId()) + 1;\n\n  const getVariables = (id: string) => [\n    { key: "POKEMON_ID", value: id },\n    { key: "BASE_URL", value: baseUrl },\n  ];\n\n  const importedPokemonList: string[] = [];\n\n  const importPokemons = async (startId: number) => {\n    const test = await tracetest.newTest(importDefinition);\n    // imports all pokemons\n    await Promise.all(\n      new Array(5).fill(0).map(async (_, index) => {\n        console.log(`\u2139 Importing pokemon ${startId + index + 1}`);\n        const run = await tracetest.runTest(test, { variables: getVariables(`${startId + index + 1}`) });\n        const updatedRun = await run.wait();\n        const pokemonId = updatedRun.outputs?.find((output) => output.name === "DATABASE_POKEMON_ID")?.value || "";\n\n        console.log(`\u2139 Adding pokemon ${pokemonId} to the list, ${updatedRun}`);\n        importedPokemonList.push(pokemonId);\n      })\n    );\n  };\n\n  const deletePokemons = async () => {\n    const test = await tracetest.newTest(deleteDefinition);\n    // deletes all pokemons\n    await Promise.all(\n      importedPokemonList.map(async (pokemonId) => {\n        console.log(`\u2139 Deleting pokemon ${pokemonId}`);\n        return tracetest.runTest(test, { variables: getVariables(pokemonId) });\n      })\n    );\n  };\n\n  await importPokemons(pokemonId);\n  console.log(await tracetest.getSummary());\n\n  await deletePokemons();\n  console.log(await tracetest.getSummary());\n};\n\nmain();\n'})}),"\n",(0,s.jsx)(t.h2,{id:"running-the-full-example",children:"Running the Full Example"}),"\n",(0,s.jsx)(t.p,{children:"To start the full setup, run the following command:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"docker compose up -d\n"})}),"\n",(0,s.jsx)(t.h2,{id:"finding-the-results",children:"Finding the Results"}),"\n",(0,s.jsxs)(t.p,{children:["The output from the TypeScript script should be visible in the log for the ",(0,s.jsx)(t.code,{children:"tracetest-client"})," service in Docker Compose. This log will show links to Tracetest for each of the test runs invoked by ",(0,s.jsx)(t.code,{children:"index.ts"}),". Click a link to launch Tracetest and view the test result."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-bash",children:"2024-01-26 10:54:44 \u2139 Importing pokemon 3\n2024-01-26 10:54:44 \u2139 Importing pokemon 4\n2024-01-26 10:54:44 \u2139 Importing pokemon 5\n2024-01-26 10:54:44 \u2139 Importing pokemon 6\n2024-01-26 10:54:44 \u2139 Importing pokemon 7\n2024-01-26 10:54:56 \u2139 Adding pokemon 1 to the list\n2024-01-26 10:54:58 \u2139 Adding pokemon 2 to the list\n2024-01-26 10:54:59 \u2139 Adding pokemon 3 to the list\n2024-01-26 10:54:59 \u2139 Adding pokemon 5 to the list\n2024-01-26 10:54:59 \u2139 Adding pokemon 4 to the list\n2024-01-26 10:54:59 \n2024-01-26 10:54:59 Successful: 5\n2024-01-26 10:54:59 Failed: 0\n2024-01-26 10:54:59 \n2024-01-26 10:54:59 [\u2714\ufe0f Typescript: Import a Pokemon] #1 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/1\n2024-01-26 10:54:59 [\u2714\ufe0f Typescript: Import a Pokemon] #2 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/2\n2024-01-26 10:54:59 [\u2714\ufe0f Typescript: Import a Pokemon] #3 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/3\n2024-01-26 10:54:59 [\u2714\ufe0f Typescript: Import a Pokemon] #4 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/4\n2024-01-26 10:54:59 [\u2714\ufe0f Typescript: Import a Pokemon] #5 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/5\n2024-01-26 10:54:59 \n2024-01-26 10:54:59 \u2139 Deleting pokemon 1\n2024-01-26 10:54:59 \u2139 Deleting pokemon 2\n2024-01-26 10:54:59 \u2139 Deleting pokemon 3\n2024-01-26 10:54:59 \u2139 Deleting pokemon 5\n2024-01-26 10:54:59 \u2139 Deleting pokemon 4\n2024-01-26 10:55:14 \n2024-01-26 10:55:14 Successful: 10\n2024-01-26 10:55:14 Failed: 0\n2024-01-26 10:55:14 \n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Import a Pokemon] #1 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/1\n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Import a Pokemon] #2 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/2\n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Import a Pokemon] #3 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/3\n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Import a Pokemon] #4 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/4\n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Import a Pokemon] #5 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/99TOHzpSR/run/5\n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Delete a Pokemon] #1 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/C2gwdktIR/run/1\n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Delete a Pokemon] #2 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/C2gwdktIR/run/2\n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Delete a Pokemon] #4 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/C2gwdktIR/run/4\n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Delete a Pokemon] #3 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/C2gwdktIR/run/3\n2024-01-26 10:55:14 [\u2714\ufe0f Typescript: Delete a Pokemon] #5 - https://app-stage.tracetest.io/organizations/ttorg_08eb62e60d1db492/environments/ttenv_70f346fe8ddba633/test/C2gwdktIR/run/5\n"})}),"\n",(0,s.jsx)(t.h2,{id:"whats-next",children:"What's Next?"}),"\n",(0,s.jsx)(t.p,{children:"After running the initial set of tests, you can click the run link for any of them, update the assertions, and run the scripts once more. This flow enables complete a trace-based TDD flow."}),"\n",(0,s.jsx)(t.p,{children:(0,s.jsx)(t.img,{alt:"assertions",src:n(31681).A+"",width:"960",height:"553"})}),"\n",(0,s.jsx)(t.h2,{id:"learn-more",children:"Learn More"}),"\n",(0,s.jsxs)(t.p,{children:["Please visit our ",(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples",children:"examples in GitHub"})," and join our ",(0,s.jsx)(t.a,{href:"https://dub.sh/tracetest-community",children:"Slack Community"})," for more info!"]})]})}function d(e={}){const{wrapper:t}={...(0,i.R)(),...e.components};return t?(0,s.jsx)(t,{...e,children:(0,s.jsx)(p,{...e})}):p(e)}},31681:(e,t,n)=>{n.d(t,{A:()=>s});const s=n.p+"assets/images/tracetest-cloud-typescript-resize-bd6c9a72c3825d10494de1a1fbf3171c.gif"},28453:(e,t,n)=>{n.d(t,{R:()=>o,x:()=>a});var s=n(96540);const i={},r=s.createContext(i);function o(e){const t=s.useContext(r);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function a(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(i):e.components||i:o(e.components),s.createElement(r.Provider,{value:t},e.children)}}}]);