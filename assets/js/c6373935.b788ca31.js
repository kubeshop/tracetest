"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[3470],{19440:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>i,default:()=>p,frontMatter:()=>a,metadata:()=>o,toc:()=>l});var r=n(74848),s=n(28453);const a={id:"running-tests-with-tracetest-graphql-pokeshop",title:"Trace-Based Tests with the Tracetest GraphQL Trigger",description:"Quickstart on how to create Trace-Based Tests with the Tracetest GraphQL Trigger",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","end-to-end testing","tracetest","graphql","trace-based-testing","trigger"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},i=void 0,o={id:"examples-tutorials/recipes/running-tests-with-tracetest-graphql-pokeshop",title:"Trace-Based Tests with the Tracetest GraphQL Trigger",description:"Quickstart on how to create Trace-Based Tests with the Tracetest GraphQL Trigger",source:"@site/docs/examples-tutorials/recipes/running-tests-with-tracetest-graphql-pokeshop.mdx",sourceDirName:"examples-tutorials/recipes",slug:"/examples-tutorials/recipes/running-tests-with-tracetest-graphql-pokeshop",permalink:"/examples-tutorials/recipes/running-tests-with-tracetest-graphql-pokeshop",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/examples-tutorials/recipes/running-tests-with-tracetest-graphql-pokeshop.mdx",tags:[],version:"current",frontMatter:{id:"running-tests-with-tracetest-graphql-pokeshop",title:"Trace-Based Tests with the Tracetest GraphQL Trigger",description:"Quickstart on how to create Trace-Based Tests with the Tracetest GraphQL Trigger",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","end-to-end testing","tracetest","graphql","trace-based-testing","trigger"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"guidesSidebar",previous:{title:"Migrating Tests from Playwright Integration to Playwright Engine",permalink:"/examples-tutorials/recipes/migrating-tests-from-playwright-integration-to-playwright-engine"},next:{title:"Synthetic Monitoring with Trace-based API Tests",permalink:"/examples-tutorials/recipes/synthetic-monitoring-trace-based-api-tests"}},c={},l=[{value:"Why is this important?",id:"why-is-this-important",level:2},{value:"Requirements",id:"requirements",level:2},{value:"Run This Example",id:"run-this-example",level:2},{value:"Project Structure",id:"project-structure",level:2},{value:"Provisioned Resources",id:"provisioned-resources",level:2},{value:"Import Pokemon Test",id:"import-pokemon-test",level:3},{value:"GraphQL",id:"graphql",level:3},{value:"GraphQL Schema",id:"graphql-schema",level:3},{value:"Jaeger Tracing Backend",id:"jaeger-tracing-backend",level:3},{value:"The Apply Script",id:"the-apply-script",level:3},{value:"The Run Script",id:"the-run-script",level:3},{value:"Setting the Environment Variables",id:"setting-the-environment-variables",level:2},{value:"Running the Full Example",id:"running-the-full-example",level:2},{value:"Finding the Results",id:"finding-the-results",level:2},{value:"What&#39;s Next?",id:"whats-next",level:2},{value:"Learn More",id:"learn-more",level:2}];function h(e){const t={a:"a",admonition:"admonition",code:"code",h2:"h2",h3:"h3",img:"img",li:"li",ol:"ol",p:"p",pre:"pre",strong:"strong",ul:"ul",...(0,s.R)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(t.admonition,{title:"Version Compatibility",type:"info",children:(0,r.jsxs)(t.p,{children:["The features described here are compatible with the ",(0,r.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/releases/tag/v1.5.2",children:"Tracetest CLI v1.5.2"})," and above."]})}),"\n",(0,r.jsx)(t.admonition,{type:"note",children:(0,r.jsx)(t.p,{children:(0,r.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples/tracetest-jaeger-graphql-pokeshop",children:"Check out the source code on GitHub here."})})}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.a,{href:"https://tracetest.io/",children:"Tracetest"})," is a testing tool based on ",(0,r.jsx)(t.a,{href:"https://opentelemetry.io/",children:"OpenTelemetry"})," that permits you to test your distributed application. It allows you to use the trace data generated by your OpenTelemetry tools to check and assert if your application has the desired behavior defined by your test definitions.\n",(0,r.jsx)(t.a,{href:"https://graphql.org/",children:"GraphQL"})," is a query language for APIs and a runtime for fulfilling those queries with your existing data. GraphQL provides a complete and understandable description of the data in your API, gives clients the power to ask for exactly what they need and nothing more, makes it easier to evolve APIs over time, and enables powerful developer tools."]}),"\n",(0,r.jsx)(t.h2,{id:"why-is-this-important",children:"Why is this important?"}),"\n",(0,r.jsx)(t.p,{children:"The Tracetest GraphQL trigger enables you to unleash the power of the trace-based testing to easily capture a full distributed trace from your OpenTelemetry instrumented GraphQL back-end system."}),"\n",(0,r.jsx)(t.p,{children:"By creating a Tracetest GraphQL test, you will be able to create trace-based assertions to be applied across the entire flow like any other Tracetest test. Not only that but it allows you to mix and match it with your existing Monitors, Test Suites and CI/CD validations."}),"\n",(0,r.jsx)(t.p,{children:"Other impactful benefits of using traces as test specs are:"}),"\n",(0,r.jsxs)(t.ul,{children:["\n",(0,r.jsx)(t.li,{children:"Faster MTTR for failing performance tests."}),"\n",(0,r.jsx)(t.li,{children:"Assert against the Mutiple Queries and Mutations at once from a single test execution."}),"\n",(0,r.jsx)(t.li,{children:"Validate functionality of other parts of your system that may be broken, even when the initial request is passing."}),"\n"]}),"\n",(0,r.jsx)(t.h2,{id:"requirements",children:"Requirements"}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.strong,{children:"Tracetest Account"}),":"]}),"\n",(0,r.jsxs)(t.ul,{children:["\n",(0,r.jsxs)(t.li,{children:["Sign up to ",(0,r.jsx)(t.a,{href:"https://app.tracetest.io",children:(0,r.jsx)(t.code,{children:"app.tracetest.io"})})," or follow the ",(0,r.jsx)(t.a,{href:"/getting-started/overview",children:"get started"})," docs."]}),"\n",(0,r.jsxs)(t.li,{children:["Have access to the environment's ",(0,r.jsx)(t.a,{href:"https://app.tracetest.io/retrieve-token",children:"agent API key"}),"."]}),"\n"]}),"\n",(0,r.jsxs)(t.p,{children:[(0,r.jsx)(t.strong,{children:"Docker"}),": Have ",(0,r.jsx)(t.a,{href:"https://docs.docker.com/get-docker/",children:"Docker"})," and ",(0,r.jsx)(t.a,{href:"https://docs.docker.com/compose/install/",children:"Docker Compose"})," installed on your machine."]}),"\n",(0,r.jsx)(t.h2,{id:"run-this-example",children:"Run This Example"}),"\n",(0,r.jsx)(t.p,{children:"The example below is provided as part of the Tracetest GitHub repo. You can download and run the example by following these steps:"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",children:"git clone https://github.com/kubeshop/tracetest\ncd tracetest/examples/tracetest-jaeger-graphql-pokeshop\n"})}),"\n",(0,r.jsx)(t.p,{children:"Follow these instructions to run the quick start:"}),"\n",(0,r.jsxs)(t.ol,{children:["\n",(0,r.jsxs)(t.li,{children:["Copy the ",(0,r.jsx)(t.code,{children:".env.template"})," file to ",(0,r.jsx)(t.code,{children:".env"}),"."]}),"\n",(0,r.jsxs)(t.li,{children:["Fill out the ",(0,r.jsx)(t.a,{href:"https://app.tracetest.io/retrieve-token",children:"TRACETEST_TOKEN and ENVIRONMENT_ID"})," details by editing your ",(0,r.jsx)(t.code,{children:".env"})," file."]}),"\n",(0,r.jsxs)(t.li,{children:["Run ",(0,r.jsx)(t.code,{children:"docker compose run tracetest-run"}),"."]}),"\n",(0,r.jsx)(t.li,{children:"Follow the links in the output to view the test results."}),"\n"]}),"\n",(0,r.jsx)(t.h2,{id:"project-structure",children:"Project Structure"}),"\n",(0,r.jsx)(t.p,{children:"The project structure for running Tracetest GraphQL tests is as follows:"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",children:".env.template\n.gitignore\n.Dockerfile\ncollector.config.yaml\ndocker-compose.yaml\n/resources\n  apply.sh\n  datastore.yaml\n  test.yaml\n  run.sh\n  scheme.graphql\n  query.graphql\n"})}),"\n",(0,r.jsxs)(t.p,{children:["The ",(0,r.jsx)(t.a,{href:"/live-examples/pokeshop/overview",children:"Pokeshop Demo App"})," is a complete example of a distributed application using different back-end and front-end services. We will be launching it and running tests against it as part of this example.\nThe ",(0,r.jsx)(t.code,{children:"docker-compose.yaml"})," file in the root directory of the quick start runs the Pokeshop Demo app, the OpenTelemetry Collector, Jaeger, and the ",(0,r.jsx)(t.a,{href:"/concepts/agent",children:"Tracetest Agent"})," setup."]}),"\n",(0,r.jsxs)(t.p,{children:["The Tracetest resource definitions and scripts are defined under the ",(0,r.jsx)(t.code,{children:"/resources"})," directory. The resources include tests and the tracing backend definition, while the scripts include the ",(0,r.jsx)(t.code,{children:"apply.sh"})," and ",(0,r.jsx)(t.code,{children:"run.sh"})," scripts to apply the resources and run the tests."]}),"\n",(0,r.jsx)(t.h2,{id:"provisioned-resources",children:"Provisioned Resources"}),"\n",(0,r.jsx)(t.p,{children:"The example provisions the following resources:"}),"\n",(0,r.jsx)(t.h3,{id:"import-pokemon-test",children:"Import Pokemon Test"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-yaml",metastring:'title="resources/test.yaml"',children:'type: Test\nspec:\n  id: re9XOxqSR\n  name: Pokeshop - Import\n  trigger:\n    type: graphql\n    graphql:\n      url: http://demo-api:8081/graphql\n      headers:\n        - key: Content-Type\n          value: application/json\n      auth:\n        apiKey: {}\n        basic: {}\n        bearer: {}\n      body:\n        query: ./query.graphql\n        variables: {}\n        operationName: ""\n      sslVerification: false\n      schema: ./schema.graphql\n  specs:\n    - name: Import Pokemon Span Exists\n      selector: span[tracetest.span.type="general" name="import pokemon"]\n      assertions:\n        - attr:tracetest.selected_spans.count = 1\n    - name: Uses Correct PokemonId\n      selector: span[tracetest.span.type="http" name="GET" http.method="GET"]\n      assertions:\n        - attr:http.url  =  "https://pokeapi.co/api/v2/pokemon/6"\n    - name: Matching db result with the Pokemon Name\n      selector: span[tracetest.span.type="database" name="create postgres.pokemon"]:first\n      assertions:\n        - attr:db.result  contains      "charizard"\n'})}),"\n",(0,r.jsx)(t.h3,{id:"graphql",children:"GraphQL"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-graphql",metastring:'title="resources/query.graphql"',children:"mutation import {\n  importPokemon(id: 6) {\n    id\n  }\n}\n"})}),"\n",(0,r.jsx)(t.h3,{id:"graphql-schema",children:"GraphQL Schema"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-graphql",metastring:'title="resources/schema.graphql"',children:'schema {\n  query: Query\n  mutation: Mutation\n}\n\ntype Pokemon {\n  id: Int\n  name: String!\n  type: String!\n  isFeatured: Boolean!\n  imageUrl: String\n}\n\n"The `Int` scalar type represents non-fractional signed whole numeric values. Int can represent values between -(2^31) and 2^31 - 1."\nscalar Int\n\n"The `String` scalar type represents textual data, represented as UTF-8 character sequences. The String type is most often used by GraphQL to represent free-form human-readable text."\nscalar String\n\n"The `Boolean` scalar type represents `true` or `false`."\nscalar Boolean\n\ntype PokemonList {\n  items: [Pokemon]\n  totalCount: Int\n}\n\ntype ImportPokemon {\n  id: Int!\n}\n\ntype Query {\n  getPokemonList(where: String, skip: Int, take: Int): PokemonList\n}\n\ntype Mutation {\n  createPokemon(name: String!, type: String!, isFeatured: Boolean!, imageUrl: String): Pokemon!\n  importPokemon(id: Int!): ImportPokemon!\n}\n// ...See more in the schema.graphql file\n'})}),"\n",(0,r.jsx)(t.h3,{id:"jaeger-tracing-backend",children:"Jaeger Tracing Backend"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-yaml",metastring:'title="resources/datastore.yaml"',children:'type: DataStore\nspec:\n  id: current\n  name: jaeger\n  type: jaeger\n  default: true\n  jaeger:\n    endpoint: jaeger:16685\n    headers:\n      "": ""\n    tls:\n      insecure: true\n'})}),"\n",(0,r.jsx)(t.h3,{id:"the-apply-script",children:"The Apply Script"}),"\n",(0,r.jsx)(t.p,{children:"The apply script configures and provisions the resources in the Tracetest environment:"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",metastring:'title="resources/apply.sh"',children:'#!/bin/sh\n\nset -e\n\nTOKEN=$TRACETEST_API_KEY\nENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID\n\napply() {\n  echo "Configuring TraceTest"\n  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID\n\n  echo "Applying Resources"\n  tracetest apply datastore -f /resources/datastore.yaml\n  tracetest apply test -f /resources/test.yaml\n}\n\napply\n'})}),"\n",(0,r.jsx)(t.h3,{id:"the-run-script",children:"The Run Script"}),"\n",(0,r.jsx)(t.p,{children:"The run script runs the test suite against the provisioned resources:"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",metastring:'title="resources/run.sh"',children:'#!/bin/sh\n\nset -e\n\nTOKEN=$TRACETEST_API_KEY\nENVIRONMENT_ID=$TRACETEST_ENVIRONMENT_ID\n\nrun() {\n  echo "Configuring Tracetest"\n  tracetest configure --token $TOKEN --environment $ENVIRONMENT_ID\n\n  echo "Running Trace-Based Tests..."\n  tracetest run test -f /resources/test.yaml\n}\n\nrun\n'})}),"\n",(0,r.jsx)(t.h2,{id:"setting-the-environment-variables",children:"Setting the Environment Variables"}),"\n",(0,r.jsxs)(t.p,{children:["Copy the ",(0,r.jsx)(t.code,{children:".env.template"})," file to ",(0,r.jsx)(t.code,{children:".env"})," and add the Tracetest API token and agent tokens to the ",(0,r.jsx)(t.code,{children:"TRACETEST_API_TOKEN"})," and ",(0,r.jsx)(t.code,{children:"TRACETEST_ENVIRONMENT_ID"})," variables."]}),"\n",(0,r.jsx)(t.h2,{id:"running-the-full-example",children:"Running the Full Example"}),"\n",(0,r.jsx)(t.p,{children:"Everything is automated for you to only run the following command:"}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",children:"docker compose run tracetest-run\n"})}),"\n",(0,r.jsxs)(t.p,{children:["This command will run the ",(0,r.jsx)(t.code,{children:"apply.sh"})," script to provision the resources and the ",(0,r.jsx)(t.code,{children:"run.sh"})," script to run the test suite."]}),"\n",(0,r.jsx)(t.h2,{id:"finding-the-results",children:"Finding the Results"}),"\n",(0,r.jsx)(t.p,{children:"The output from the Tracetest Engine script should be visible in the console log after running the test command."}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{className:"language-bash",metastring:"title=Output",children:"WARN[0000] /Users/oscar/Documents/kubeshop/t/examples/tracetest-jaeger-graphql-pokeshop/docker-compose.yaml: `version` is obsolete\n[+] Running 2/2\n \u2714 demo-api Pulled                                                                                                                                                0.9s\n \u2714 demo-worker Pulled                                                                                                                                             0.9s\n[+] Creating 10/9\n \u2714 Network tracetest-jaeger-graphql-pokeshop_default              Created                                                                                         0.0s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-postgres-1         Created                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-cache-1            Created                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-tracetest-agent-1  Created                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-queue-1            Created                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-jaeger-1           Created                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-demo-worker-1      Created                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-otel-collector-1   Created                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-demo-api-1         Created                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-tracetest-apply-1  Created                                                                                         0.0s\n[+] Running 9/9\n \u2714 Container tracetest-jaeger-graphql-pokeshop-cache-1            Healthy                                                                                        10.5s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-tracetest-agent-1  Started                                                                                         0.3s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-jaeger-1           Healthy                                                                                         1.9s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-postgres-1         Healthy                                                                                        10.5s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-queue-1            Healthy                                                                                        10.5s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-otel-collector-1   Started                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-demo-worker-1      Started                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-demo-api-1         Started                                                                                         0.1s\n \u2714 Container tracetest-jaeger-graphql-pokeshop-tracetest-apply-1  Started                                                                                         0.1s\n[+] Running 2/2\n \u2714 demo-api Pulled                                                                                                                                                0.8s\n \u2714 demo-worker Pulled                                                                                                                                             0.8s\nConfiguring Tracetest\n SUCCESS  Successfully configured Tracetest CLI\nRunning Trace-Based Tests...\n\u2714 RunGroup: #US5klbqSR (https://app-stage.tracetest.io/organizations/ttorg_c71a6b53c3709e95/environments/ttenv_bcf29b43f06a12dc/run/US5klbqSR)\n Summary: 1 passed, 0 failed, 0 pending\n  \u2714 Pokeshop - Import (https://app-stage.tracetest.io/organizations/ttorg_c71a6b53c3709e95/environments/ttenv_bcf29b43f06a12dc/test/re9XOxqSR/run/11/test) - trace id: 6facf84ee23757eda97930c16fd1d8f9\n\t\u2714 Import Pokemon Span Exists\n\t\u2714 Uses Correct PokemonId\n\t\u2714 Matching db result with the Pokemon Name\n"})}),"\n",(0,r.jsx)(t.h2,{id:"whats-next",children:"What's Next?"}),"\n",(0,r.jsx)(t.p,{children:"After running the test, you can click the run link, update the assertions, and run the scripts once more. This flow enables complete a trace-based TDD flow."}),"\n",(0,r.jsx)(t.p,{children:(0,r.jsx)(t.img,{alt:"assertions",src:n(76667).A+"",width:"1650",height:"1255"})}),"\n",(0,r.jsx)(t.h2,{id:"learn-more",children:"Learn More"}),"\n",(0,r.jsxs)(t.p,{children:["Please visit our ",(0,r.jsx)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples",children:"examples in GitHub"})," and join our ",(0,r.jsx)(t.a,{href:"https://dub.sh/tracetest-community",children:"Slack Community"})," for more info!"]})]})}function p(e={}){const{wrapper:t}={...(0,s.R)(),...e.components};return t?(0,r.jsx)(t,{...e,children:(0,r.jsx)(h,{...e})}):h(e)}},76667:(e,t,n)=>{n.d(t,{A:()=>r});const r=n.p+"assets/images/playwright-engine-fd1b71b240a697bdcb171a6c80c9f4a2.gif"},28453:(e,t,n)=>{n.d(t,{R:()=>i,x:()=>o});var r=n(96540);const s={},a=r.createContext(s);function i(e){const t=r.useContext(a);return r.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function o(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(s):e.components||s:i(e.components),r.createElement(a.Provider,{value:t},e.children)}}}]);