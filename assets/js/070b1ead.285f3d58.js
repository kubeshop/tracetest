"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[2183],{81661:(e,s,t)=>{t.r(s),t.d(s,{assets:()=>o,contentTitle:()=>i,default:()=>p,frontMatter:()=>a,metadata:()=>l,toc:()=>d});var n=t(74848),r=t(28453);const a={id:"no-api-key-leak",title:"no-api-key-leak",description:"Disallow leaked API keys for HTTP spans | The Tracetest Analyzer analyzes OpenTelemetry traces",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},i=void 0,l={id:"analyzer/rules/no-api-key-leak",title:"no-api-key-leak",description:"Disallow leaked API keys for HTTP spans | The Tracetest Analyzer analyzes OpenTelemetry traces",source:"@site/docs/analyzer/rules/no-api-key-leak.mdx",sourceDirName:"analyzer/rules",slug:"/analyzer/rules/no-api-key-leak",permalink:"/analyzer/rules/no-api-key-leak",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/analyzer/rules/no-api-key-leak.mdx",tags:[],version:"current",frontMatter:{id:"no-api-key-leak",title:"no-api-key-leak",description:"Disallow leaked API keys for HTTP spans | The Tracetest Analyzer analyzes OpenTelemetry traces",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"tutorialSidebar",previous:{title:"secure-https-protocol",permalink:"/analyzer/rules/secure-https-protocol"},next:{title:"Common Problems",permalink:"/analyzer/plugins/common-problems"}},o={},d=[{value:"Rule Details",id:"rule-details",level:2},{value:"HTTP spans:",id:"http-spans",level:3},{value:"Options",id:"options",level:2},{value:"When Not To Use It",id:"when-not-to-use-it",level:2}];function c(e){const s={code:"code",h2:"h2",h3:"h3",li:"li",p:"p",pre:"pre",ul:"ul",...(0,r.R)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(s.p,{children:"Disallow leaked API keys for HTTP spans."}),"\n",(0,n.jsx)(s.h2,{id:"rule-details",children:"Rule Details"}),"\n",(0,n.jsx)(s.p,{children:"This rule disallows the recording of API keys in HTTP spans."}),"\n",(0,n.jsx)(s.h3,{id:"http-spans",children:"HTTP spans:"}),"\n",(0,n.jsx)(s.p,{children:"The following attributes are evaluated:"}),"\n",(0,n.jsx)(s.pre,{children:(0,n.jsx)(s.code,{children:"- http.response.header.authorization\n- http.response.header.x-api-key\n- http.request.header.authorization\n- http.request.header.x-api-key\n"})}),"\n",(0,n.jsx)(s.h2,{id:"options",children:"Options"}),"\n",(0,n.jsx)(s.p,{children:"This rule has the following options:"}),"\n",(0,n.jsxs)(s.ul,{children:["\n",(0,n.jsxs)(s.li,{children:[(0,n.jsx)(s.code,{children:'"error"'})," requires no leaked API keys for HTTP spans"]}),"\n",(0,n.jsxs)(s.li,{children:[(0,n.jsx)(s.code,{children:'"disabled"'})," disables the no leaked API keys verification for HTTP spans"]}),"\n",(0,n.jsxs)(s.li,{children:[(0,n.jsx)(s.code,{children:'"warning"'})," verifies no leaked API keys for HTTPS spans but does not impact the analyzer score"]}),"\n"]}),"\n",(0,n.jsx)(s.h2,{id:"when-not-to-use-it",children:"When Not To Use It"}),"\n",(0,n.jsx)(s.p,{children:"If you intentionally record API keys for HTTP spans then you can disable this rule."})]})}function p(e={}){const{wrapper:s}={...(0,r.R)(),...e.components};return s?(0,n.jsx)(s,{...e,children:(0,n.jsx)(c,{...e})}):c(e)}},28453:(e,s,t)=>{t.d(s,{R:()=>i,x:()=>l});var n=t(96540);const r={},a=n.createContext(r);function i(e){const s=n.useContext(a);return n.useMemo((function(){return"function"==typeof e?e(s):{...s,...e}}),[s,e])}function l(e){let s;return s=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:i(e.components),n.createElement(a.Provider,{value:s},e.children)}}}]);