"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[2651],{71813:(e,t,s)=>{s.r(t),s.d(t,{assets:()=>i,contentTitle:()=>l,default:()=>u,frontMatter:()=>n,metadata:()=>c,toc:()=>a});var r=s(74848),o=s(28453);const n={id:"secure-https-protocol",title:"secure-https-protocol",description:"Enforce usage of secure protocol for HTTP server spans | The Tracetest Analyzer analyzes OpenTelemetry traces",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},l=void 0,c={id:"analyzer/rules/secure-https-protocol",title:"secure-https-protocol",description:"Enforce usage of secure protocol for HTTP server spans | The Tracetest Analyzer analyzes OpenTelemetry traces",source:"@site/docs/analyzer/rules/secure-https-protocol.mdx",sourceDirName:"analyzer/rules",slug:"/analyzer/rules/secure-https-protocol",permalink:"/analyzer/rules/secure-https-protocol",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/analyzer/rules/secure-https-protocol.mdx",tags:[],version:"current",frontMatter:{id:"secure-https-protocol",title:"secure-https-protocol",description:"Enforce usage of secure protocol for HTTP server spans | The Tracetest Analyzer analyzes OpenTelemetry traces",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"tutorialSidebar",previous:{title:"Security",permalink:"/analyzer/plugins/security"},next:{title:"no-api-key-leak",permalink:"/analyzer/rules/no-api-key-leak"}},i={},a=[{value:"Rule Details",id:"rule-details",level:2},{value:"HTTP spans:",id:"http-spans",level:3},{value:"Options",id:"options",level:2},{value:"When Not To Use It",id:"when-not-to-use-it",level:2}];function p(e){const t={code:"code",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,o.R)(),...e.components};return(0,r.jsxs)(r.Fragment,{children:[(0,r.jsx)(t.p,{children:"Enforce usage of secure protocol for HTTP server spans."}),"\n",(0,r.jsx)(t.h2,{id:"rule-details",children:"Rule Details"}),"\n",(0,r.jsxs)(t.p,{children:["This rule enforces usage of a secure protocol for an HTTP server span. The URI scheme that identifies the used protocol should be ",(0,r.jsx)(t.code,{children:'"https"'}),"."]}),"\n",(0,r.jsx)(t.h3,{id:"http-spans",children:"HTTP spans:"}),"\n",(0,r.jsxs)(t.p,{children:["If span kind is ",(0,r.jsx)(t.code,{children:'"server"'}),", the following attributes are evaluated:"]}),"\n",(0,r.jsx)(t.pre,{children:(0,r.jsx)(t.code,{children:'- http.scheme = "https"\n- http.url = "https"\n'})}),"\n",(0,r.jsx)(t.h2,{id:"options",children:"Options"}),"\n",(0,r.jsx)(t.p,{children:"This rule has the following options:"}),"\n",(0,r.jsxs)(t.ul,{children:["\n",(0,r.jsxs)(t.li,{children:[(0,r.jsx)(t.code,{children:'"error"'})," requires secure protocol for HTTP server spans"]}),"\n",(0,r.jsxs)(t.li,{children:[(0,r.jsx)(t.code,{children:'"disabled"'})," disables the secure protocol verification for HTTP server spans"]}),"\n",(0,r.jsxs)(t.li,{children:[(0,r.jsx)(t.code,{children:'"warning"'})," verifies secure protocol for HTTPS server spans but does not impact the analyzer score"]}),"\n"]}),"\n",(0,r.jsx)(t.h2,{id:"when-not-to-use-it",children:"When Not To Use It"}),"\n",(0,r.jsx)(t.p,{children:"If you intentionally use non secure protocol for HTTP server spans then you can disable this rule."})]})}function u(e={}){const{wrapper:t}={...(0,o.R)(),...e.components};return t?(0,r.jsx)(t,{...e,children:(0,r.jsx)(p,{...e})}):p(e)}},28453:(e,t,s)=>{s.d(t,{R:()=>l,x:()=>c});var r=s(96540);const o={},n=r.createContext(o);function l(e){const t=r.useContext(n);return r.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function c(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(o):e.components||o:l(e.components),r.createElement(n.Provider,{value:t},e.children)}}}]);