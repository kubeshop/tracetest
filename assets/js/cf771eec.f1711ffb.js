"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[3808],{40335:(e,t,s)=>{s.r(t),s.d(t,{assets:()=>o,contentTitle:()=>r,default:()=>h,frontMatter:()=>a,metadata:()=>d,toc:()=>c});var i=s(74848),n=s(28453);const a={id:"ad-hoc-testing",title:"Ad-hoc Testing",description:"Use Tracetest to enable ad-hoc testing by utilizing variable sets and undefined variables.",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},r=void 0,d={id:"concepts/ad-hoc-testing",title:"Ad-hoc Testing",description:"Use Tracetest to enable ad-hoc testing by utilizing variable sets and undefined variables.",source:"@site/docs/concepts/ad-hoc-testing.mdx",sourceDirName:"concepts",slug:"/concepts/ad-hoc-testing",permalink:"/concepts/ad-hoc-testing",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/concepts/ad-hoc-testing.mdx",tags:[],version:"current",frontMatter:{id:"ad-hoc-testing",title:"Ad-hoc Testing",description:"Use Tracetest to enable ad-hoc testing by utilizing variable sets and undefined variables.",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"tutorialSidebar",previous:{title:"Variable Sets",permalink:"/concepts/variable-sets"},next:{title:"Trace Analyzer Concepts",permalink:"/analyzer/concepts"}},o={},c=[{value:"Undefined Variables Use Cases",id:"undefined-variables-use-cases",level:2},{value:"Supply Variable Value at Runtime",id:"supply-variable-value-at-runtime",level:3},{value:"Supply Variable Value from a Previous Test",id:"supply-variable-value-from-a-previous-test",level:3},{value:"Undefined Variables Test Suite with Multiple Tests Example",id:"undefined-variables-test-suite-with-multiple-tests-example",level:2}];function l(e){const t={a:"a",admonition:"admonition",code:"code",h2:"h2",h3:"h3",img:"img",li:"li",ol:"ol",p:"p",strong:"strong",...(0,n.R)(),...e.components};return(0,i.jsxs)(i.Fragment,{children:[(0,i.jsx)(t.p,{children:"This page showcases use-cases for undefined variables and how to enable ad-hoc testing by utilizing variable sets and undefined variables."}),"\n",(0,i.jsx)(t.admonition,{type:"tip",children:(0,i.jsx)(t.p,{children:(0,i.jsxs)(t.a,{href:"/web-ui/undefined-variables",children:["Check out how to configure ad-hoc testing with undefined variables with the ",(0,i.jsx)(t.strong,{children:"Web UI"})," here."]})})}),"\n",(0,i.jsx)(t.admonition,{type:"tip",children:(0,i.jsx)(t.p,{children:(0,i.jsxs)(t.a,{href:"/cli/undefined-variables",children:["Check out how to configure ad-hoc testing with undefined variables with the ",(0,i.jsx)(t.strong,{children:"CLI"})," here."]})})}),"\n",(0,i.jsx)(t.h2,{id:"undefined-variables-use-cases",children:"Undefined Variables Use Cases"}),"\n",(0,i.jsx)(t.h3,{id:"supply-variable-value-at-runtime",children:"Supply Variable Value at Runtime"}),"\n",(0,i.jsx)(t.p,{children:"A user wants a Test or Test Suite they can run on a particular user, order id, etc. that is configurable at run time. This makes running an adhoc test in an environment, even production, very easy and convenient. In this case, the user references the variable, but doesn't add it to the environment. Each time they run the Test or Test Suite, they will be prompted for the unspecified variables."}),"\n",(0,i.jsx)(t.h3,{id:"supply-variable-value-from-a-previous-test",children:"Supply Variable Value from a Previous Test"}),"\n",(0,i.jsx)(t.p,{children:"A user wants to define 3 tests as part of a Test Suite. The first test has an output variable and this output is used by the second test. They define the first test. They then define the second test and reference the variable value that is output from the first test."}),"\n",(0,i.jsx)(t.p,{children:"In Tracetest, undefined variables can be used in both the UI and CLI."}),"\n",(0,i.jsx)(t.h2,{id:"undefined-variables-test-suite-with-multiple-tests-example",children:"Undefined Variables Test Suite with Multiple Tests Example"}),"\n",(0,i.jsxs)(t.ol,{children:["\n",(0,i.jsx)(t.li,{children:"Create an HTTP Pokemon list test that uses variables for hostname and the SKIP query parameter:"}),"\n"]}),"\n",(0,i.jsx)(t.p,{children:(0,i.jsx)(t.img,{alt:"Create Pokemon List",src:s(65168).A+"",width:"2874",height:"1436"})}),"\n",(0,i.jsxs)(t.ol,{children:["\n",(0,i.jsxs)(t.li,{children:["Within the test, create test spec assertions that use variables for comparators, something like: ",(0,i.jsx)(t.code,{children:'http.status_code = "${var:STATUS_CODE}"'}),":"]}),"\n"]}),"\n",(0,i.jsx)(t.p,{children:(0,i.jsx)(t.img,{alt:"Create Test Spec Assertionsl",src:s(47967).A+"",width:"2874",height:"1474"})}),"\n",(0,i.jsxs)(t.ol,{children:["\n",(0,i.jsx)(t.li,{children:"Create a GRPC Pokemon add test that uses variables for hostname and Pokemon name:"}),"\n"]}),"\n",(0,i.jsx)(t.p,{children:(0,i.jsx)(t.img,{alt:"Create GRPC",src:s(9705).A+"",width:"2854",height:"1466"})}),"\n",(0,i.jsxs)(t.ol,{start:"4",children:["\n",(0,i.jsx)(t.li,{children:"Create an output from this test for the SKIP variable that could come from anywhere in the trace:"}),"\n"]}),"\n",(0,i.jsx)(t.p,{children:(0,i.jsx)(t.img,{alt:"Test Output",src:s(63717).A+"",width:"2874",height:"1470"})}),"\n",(0,i.jsxs)(t.ol,{start:"5",children:["\n",(0,i.jsx)(t.li,{children:"Now, you can create a Test Suite with the two tests - first, add the list test, then the add test, and then the list test again:"}),"\n"]}),"\n",(0,i.jsx)(t.p,{children:(0,i.jsx)(t.img,{src:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1710422676/docs/create-testsuite-min_jawmao.png",alt:"Create Test Suite"})}),"\n",(0,i.jsxs)(t.ol,{start:"6",children:["\n",(0,i.jsx)(t.li,{children:"From here you can input the values for the undefined variables and complete your trace:"}),"\n"]}),"\n",(0,i.jsx)(t.p,{children:(0,i.jsx)(t.img,{alt:"Input Values",src:s(77161).A+"",width:"2854",height:"1462"})})]})}function h(e={}){const{wrapper:t}={...(0,n.R)(),...e.components};return t?(0,i.jsx)(t,{...e,children:(0,i.jsx)(l,{...e})}):l(e)}},9705:(e,t,s)=>{s.d(t,{A:()=>i});const i=s.p+"assets/images/create-grpc-d0ac8db6677436dd601053bdc1ddc7e0.png"},47967:(e,t,s)=>{s.d(t,{A:()=>i});const i=s.p+"assets/images/create-test-spec-assertions-3d718497ce13b7260b8a094be21218cf.png"},77161:(e,t,s)=>{s.d(t,{A:()=>i});const i=s.p+"assets/images/input-values-ed1d74e1c31d1194dd212297e3b3e6b6.png"},65168:(e,t,s)=>{s.d(t,{A:()=>i});const i=s.p+"assets/images/pokeshop-list-77a8013a017df9404efca27a024a53fe.png"},63717:(e,t,s)=>{s.d(t,{A:()=>i});const i=s.p+"assets/images/test-output-4728ba75afad2c8ed8cffc8c925c3ce2.png"},28453:(e,t,s)=>{s.d(t,{R:()=>r,x:()=>d});var i=s(96540);const n={},a=i.createContext(n);function r(e){const t=i.useContext(a);return i.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function d(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(n):e.components||n:r(e.components),i.createElement(a.Provider,{value:t},e.children)}}}]);