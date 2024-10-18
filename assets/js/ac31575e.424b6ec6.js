"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[17],{63875:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>i,default:()=>h,frontMatter:()=>o,metadata:()=>a,toc:()=>d});var s=n(74848),r=n(28453);const o={},i="Trace Mode",a={id:"web-ui/trace-mode",title:"Trace Mode",description:"Trace Mode enables you to verify trace ingestion is configured correctly, and is also as a starting point to create trace-based tests.",source:"@site/docs/web-ui/trace-mode.md",sourceDirName:"web-ui",slug:"/web-ui/trace-mode",permalink:"/web-ui/trace-mode",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/web-ui/trace-mode.md",tags:[],version:"current",frontMatter:{},sidebar:"tutorialSidebar",previous:{title:"Configuring Trace Ingestion",permalink:"/web-ui/creating-data-stores"},next:{title:"Creating Tests",permalink:"/web-ui/creating-tests"}},c={},d=[];function l(e){const t={a:"a",admonition:"admonition",code:"code",h1:"h1",img:"img",li:"li",p:"p",strong:"strong",ul:"ul",...(0,r.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(t.h1,{id:"trace-mode",children:"Trace Mode"}),"\n",(0,s.jsx)(t.p,{children:"Trace Mode enables you to verify trace ingestion is configured correctly, and is also as a starting point to create trace-based tests."}),"\n",(0,s.jsx)(t.p,{children:"You can use it to:"}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsx)(t.li,{children:"View traces coming in to your Tracetest account."}),"\n",(0,s.jsx)(t.li,{children:"Create tests from certain trace spans. Tracetest figures out how to help you trigger tests based on span attributes and metadata."}),"\n",(0,s.jsx)(t.li,{children:"Create tests from trace IDs."}),"\n"]}),"\n",(0,s.jsx)(t.admonition,{type:"note",children:(0,s.jsx)(t.p,{children:"Traces will be deleted after 4 days."})}),"\n",(0,s.jsx)(t.p,{children:"The steps to use Tracetest Trace Mode are:"}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"1."})," Pull the latest version of the Pokeshop repo master branch ",(0,s.jsx)(t.a,{href:"https://github.com/kubeshop/pokeshop",children:"here"}),"."]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"2."})," Make sure Docker is running."]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"3."})," In the Docker Desktop, search for and make sure to delete any previous agent image."]}),"\n",(0,s.jsx)(t.p,{children:(0,s.jsx)(t.img,{alt:"Delete Previous Agent",src:n(14789).A+"",width:"1844",height:"446"})}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"4."})," Create the .env file in the pokeshop root folder from the template and add\n",(0,s.jsx)(t.code,{children:"POKESHOP_DEMO_URL=http://localhost:8081 TRACETEST_AGENT_API_KEY=<your-agent-api-key> TRACETEST_ENVIRONMENT_ID=<your-environment-id> TRACETEST_TRACE_MODE=true"}),"."]}),"\n",(0,s.jsxs)(t.p,{children:["The agent API key and Environment ID can be found ",(0,s.jsx)(t.a,{href:"https://app.tracetest.io/retrieve-token",children:"here"}),"."]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"5."})," From the ",(0,s.jsx)(t.code,{children:"pokeshop"})," root folder run ",(0,s.jsx)(t.code,{children:"docker compose -f docker-compose.yml -f docker-compose.e2e.yml up -d"}),"."]}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"6."})," From the Tracetest app, in ",(0,s.jsx)(t.strong,{children:"Settings"}),", go to the ",(0,s.jsx)(t.strong,{children:"Trace Ingestion"})," configuration tab and select ",(0,s.jsx)(t.strong,{children:"Open Telemetry"}),"."]}),"\n",(0,s.jsx)(t.p,{children:(0,s.jsx)(t.img,{alt:"Trace Ingestion",src:n(72770).A+"",width:"2874",height:"1572"})}),"\n",(0,s.jsxs)(t.p,{children:["You will see the ",(0,s.jsx)(t.strong,{children:"Open Telemetry"})," details and click ",(0,s.jsx)(t.strong,{children:"Save"}),":"]}),"\n",(0,s.jsx)(t.p,{children:(0,s.jsx)(t.img,{alt:"Trace Ingestion Save",src:n(99430).A+"",width:"2874",height:"1574"})}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"7."})," Open your environment dashboard and look at the Traces' landing page."]}),"\n",(0,s.jsx)(t.p,{children:(0,s.jsx)(t.img,{alt:"Trace Landing Page",src:n(42227).A+"",width:"2512",height:"1570"})}),"\n",(0,s.jsxs)(t.p,{children:[(0,s.jsx)(t.strong,{children:"8."})," Play around with the Pokeshop UI at ",(0,s.jsx)(t.a,{href:"http://localhost:8081",children:"http://localhost:8081"}),". You'll see traces appearing in the landing page."]}),"\n",(0,s.jsx)(t.p,{children:(0,s.jsx)(t.img,{alt:"Trace Details",src:n(15090).A+"",width:"2508",height:"1572"})})]})}function h(e={}){const{wrapper:t}={...(0,r.R)(),...e.components};return t?(0,s.jsx)(t,{...e,children:(0,s.jsx)(l,{...e})}):l(e)}},14789:(e,t,n)=>{n.d(t,{A:()=>s});const s=n.p+"assets/images/delete-previous-agent-d7c80a525fd764fb0ecad41fb4bafd90.png"},15090:(e,t,n)=>{n.d(t,{A:()=>s});const s=n.p+"assets/images/trace-details-e9ea9d75fdded6016620d82150596d4b.png"},99430:(e,t,n)=>{n.d(t,{A:()=>s});const s=n.p+"assets/images/trace-ingestion-save-f2b8ad34f806a3833e14ae323d28aa74.png"},72770:(e,t,n)=>{n.d(t,{A:()=>s});const s=n.p+"assets/images/trace-ingestion-6af23e4b35c749620b4e3dc489f7d496.png"},42227:(e,t,n)=>{n.d(t,{A:()=>s});const s=n.p+"assets/images/traces-list-b09a52df71e52990edaddb8f8b5ff813.png"},28453:(e,t,n)=>{n.d(t,{R:()=>i,x:()=>a});var s=n(96540);const r={},o=s.createContext(r);function i(e){const t=s.useContext(o);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function a(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:i(e.components),s.createElement(o.Provider,{value:t},e.children)}}}]);