"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[2362],{29116:(e,t,o)=>{o.r(t),o.d(t,{assets:()=>i,contentTitle:()=>s,default:()=>h,frontMatter:()=>r,metadata:()=>a,toc:()=>d});var n=o(74848),c=o(28453);const r={id:"honeycomb",title:"Honeycomb",description:"Use Honeycomb as the trace data store for Tracetest. Configure the OpenTelemetry Collector to receive traces and forward them to both Tracetest and Honeycomb.",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},s=void 0,a={id:"configuration/connecting-to-data-stores/honeycomb",title:"Honeycomb",description:"Use Honeycomb as the trace data store for Tracetest. Configure the OpenTelemetry Collector to receive traces and forward them to both Tracetest and Honeycomb.",source:"@site/docs/configuration/connecting-to-data-stores/honeycomb.mdx",sourceDirName:"configuration/connecting-to-data-stores",slug:"/configuration/connecting-to-data-stores/honeycomb",permalink:"/configuration/connecting-to-data-stores/honeycomb",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/configuration/connecting-to-data-stores/honeycomb.mdx",tags:[],version:"current",frontMatter:{id:"honeycomb",title:"Honeycomb",description:"Use Honeycomb as the trace data store for Tracetest. Configure the OpenTelemetry Collector to receive traces and forward them to both Tracetest and Honeycomb.",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"tutorialSidebar",previous:{title:"Grafana Tempo",permalink:"/configuration/connecting-to-data-stores/tempo"},next:{title:"Instana",permalink:"/configuration/connecting-to-data-stores/instana"}},i={},d=[{value:"Configuring OpenTelemetry Collector to Send Traces to both Honeycomb and Tracetest",id:"configuring-opentelemetry-collector-to-send-traces-to-both-honeycomb-and-tracetest",level:2},{value:"Configure Tracetest to Use Honeycomb as a Trace Data Store",id:"configure-tracetest-to-use-honeycomb-as-a-trace-data-store",level:2},{value:"Connect Tracetest to Honeycomb with the Web UI",id:"connect-tracetest-to-honeycomb-with-the-web-ui",level:2},{value:"Connect Tracetest to Honeycomb with the CLI",id:"connect-tracetest-to-honeycomb-with-the-cli",level:2}];function l(e){const t={a:"a",admonition:"admonition",code:"code",h2:"h2",img:"img",li:"li",p:"p",pre:"pre",ul:"ul",...(0,c.R)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsxs)(t.p,{children:["If you want to use ",(0,n.jsx)(t.a,{href:"https://honeycomb.io/",children:"Honeycomb"})," as the trace data store, you'll configure the OpenTelemetry Collector to receive traces from your system and then send them to both Tracetest and Honeycomb. And, you don't have to change your existing pipelines to do so."]}),"\n",(0,n.jsx)(t.admonition,{type:"tip",children:(0,n.jsxs)(t.p,{children:["Examples of configuring Tracetest with Honeycomb can be found in the ",(0,n.jsxs)(t.a,{href:"https://github.com/kubeshop/tracetest/tree/main/examples",children:[(0,n.jsx)(t.code,{children:"examples"})," folder of the Tracetest GitHub repo"]}),"."]})}),"\n",(0,n.jsx)(t.h2,{id:"configuring-opentelemetry-collector-to-send-traces-to-both-honeycomb-and-tracetest",children:"Configuring OpenTelemetry Collector to Send Traces to both Honeycomb and Tracetest"}),"\n",(0,n.jsx)(t.p,{children:"In your OpenTelemetry Collector config file:"}),"\n",(0,n.jsxs)(t.ul,{children:["\n",(0,n.jsxs)(t.li,{children:["Set the ",(0,n.jsx)(t.code,{children:"exporter"})," to ",(0,n.jsx)(t.code,{children:"otlp/tracetest"})]}),"\n",(0,n.jsxs)(t.li,{children:["Set the ",(0,n.jsx)(t.code,{children:"endpoint"})," to your Tracetest instance on port ",(0,n.jsx)(t.code,{children:"4317"})]}),"\n"]}),"\n",(0,n.jsx)(t.admonition,{type:"tip",children:(0,n.jsxs)(t.p,{children:["If you are running Tracetest with Docker, and Tracetest's service name is ",(0,n.jsx)(t.code,{children:"tracetest"}),", then the endpoint might look like this ",(0,n.jsx)(t.code,{children:"http://tracetest:4317"})]})}),"\n",(0,n.jsx)(t.p,{children:"Additionally, add another config:"}),"\n",(0,n.jsxs)(t.ul,{children:["\n",(0,n.jsxs)(t.li,{children:["Set the ",(0,n.jsx)(t.code,{children:"exporter"})," to ",(0,n.jsx)(t.code,{children:"otlp/honeycomb"})]}),"\n",(0,n.jsxs)(t.li,{children:["Set the ",(0,n.jsx)(t.code,{children:"endpoint"})," pointing to the Honeycomb API and using Honeycomb API KEY"]}),"\n"]}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-yaml",children:'# collector.config.yaml\n\n# If you already have receivers declared, you can just ignore\n# this one and still use yours instead.\nreceivers:\n  otlp:\n    protocols:\n      grpc:\n      http:\n\nprocessors:\n  batch:\n    timeout: 100ms\n\nexporters:\n  logging:\n    logLevel: debug\n  # OTLP for Tracetest\n  otlp/tracetest:\n    endpoint: tracetest:4317 # Send traces to Tracetest. Read more in docs here:  https://docs.tracetest.io/configuration/connecting-to-data-stores/opentelemetry-collector\n    tls:\n      insecure: true\n  # OTLP for Honeycomb\n  otlp/honeycomb:\n    endpoint: "api.honeycomb.io:443"\n    headers:\n      "x-honeycomb-team": "YOUR_API_KEY"\n      # Read more in docs here: https://docs.honeycomb.io/getting-data-in/otel-collector/\n\nservice:\n  pipelines:\n    traces/tracetest:\n      receivers: [otlp]\n      processors: [batch]\n      exporters: [otlp/tracetest]\n    traces/honeycomb:\n      receivers: [otlp]\n      processors: [batch]\n      exporters: [logging, otlp/honeycomb]\n'})}),"\n",(0,n.jsx)(t.h2,{id:"configure-tracetest-to-use-honeycomb-as-a-trace-data-store",children:"Configure Tracetest to Use Honeycomb as a Trace Data Store"}),"\n",(0,n.jsxs)(t.p,{children:["Configure your Tracetest instance to expose an ",(0,n.jsx)(t.code,{children:"otlp"})," endpoint to make it aware it will receive traces from the OpenTelemetry Collector. This will expose Tracetest's trace receiver on port ",(0,n.jsx)(t.code,{children:"4317"}),"."]}),"\n",(0,n.jsx)(t.h2,{id:"connect-tracetest-to-honeycomb-with-the-web-ui",children:"Connect Tracetest to Honeycomb with the Web UI"}),"\n",(0,n.jsx)(t.p,{children:"In the Web UI, (1) open Settings, and, on the (2) Trace Ingestion tab, select (3) Honeycomb."}),"\n",(0,n.jsx)(t.p,{children:(0,n.jsx)(t.img,{alt:"Trace Ingestion Settings",src:o(5451).A+"",width:"3326",height:"1808"})}),"\n",(0,n.jsx)(t.h2,{id:"connect-tracetest-to-honeycomb-with-the-cli",children:"Connect Tracetest to Honeycomb with the CLI"}),"\n",(0,n.jsx)(t.p,{children:"Or, if you prefer using the CLI, you can use this file config."}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-yaml",children:"type: DataStore\nspec:\n  name: Honeycomb pipeline\n  type: honeycomb\n  default: true\n"})}),"\n",(0,n.jsx)(t.p,{children:"Proceed to run this command in the terminal and specify the file above."}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-bash",children:"tracetest apply datastore -f my/data-store/file/location.yaml\n"})}),"\n",(0,n.jsx)(t.admonition,{type:"tip",children:(0,n.jsxs)(t.p,{children:["To learn more, ",(0,n.jsx)(t.a,{href:"/examples-tutorials/recipes/running-tracetest-with-honeycomb",children:"read the recipe on running a sample app with Honeycomb and Tracetest"}),"."]})})]})}function h(e={}){const{wrapper:t}={...(0,c.R)(),...e.components};return t?(0,n.jsx)(t,{...e,children:(0,n.jsx)(l,{...e})}):l(e)}},5451:(e,t,o)=>{o.d(t,{A:()=>n});const n=o.p+"assets/images/app.tracetest.io_organizations_at4CxvjIg_environments_ttenv_172de56e3dcbba9b_settings_tab=dataStore_honeycomb-3812e51d52b2dff4bf32d63ea7d15a52.png"},28453:(e,t,o)=>{o.d(t,{R:()=>s,x:()=>a});var n=o(96540);const c={},r=n.createContext(c);function s(e){const t=n.useContext(r);return n.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function a(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(c):e.components||c:s(e.components),n.createElement(r.Provider,{value:t},e.children)}}}]);