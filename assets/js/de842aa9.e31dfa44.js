"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[4043],{64477:(e,s,n)=>{n.r(s),n.d(s,{assets:()=>d,contentTitle:()=>o,default:()=>u,frontMatter:()=>i,metadata:()=>a,toc:()=>l});var t=n(74848),r=n(28453);const i={id:"prefer-dns",title:"prefer-dns",description:"Enforce usage of DNS instead of IP addresses | The Tracetest Analyzer analyzes OpenTelemetry traces",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},o=void 0,a={id:"analyzer/rules/prefer-dns",title:"prefer-dns",description:"Enforce usage of DNS instead of IP addresses | The Tracetest Analyzer analyzes OpenTelemetry traces",source:"@site/docs/analyzer/rules/prefer-dns.mdx",sourceDirName:"analyzer/rules",slug:"/analyzer/rules/prefer-dns",permalink:"/analyzer/rules/prefer-dns",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/analyzer/rules/prefer-dns.mdx",tags:[],version:"current",frontMatter:{id:"prefer-dns",title:"prefer-dns",description:"Enforce usage of DNS instead of IP addresses | The Tracetest Analyzer analyzes OpenTelemetry traces",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"tutorialSidebar",previous:{title:"Common Problems",permalink:"/analyzer/plugins/common-problems"},next:{title:"Environment Automation",permalink:"/concepts/environment-automation"}},d={},l=[{value:"Rule Details",id:"rule-details",level:2},{value:"Options",id:"options",level:2},{value:"When Not To Use It",id:"when-not-to-use-it",level:2}];function c(e){const s={code:"code",h2:"h2",li:"li",p:"p",pre:"pre",ul:"ul",...(0,r.R)(),...e.components};return(0,t.jsxs)(t.Fragment,{children:[(0,t.jsx)(s.p,{children:"Enforce usage of DNS instead of IP addresses."}),"\n",(0,t.jsx)(s.h2,{id:"rule-details",children:"Rule Details"}),"\n",(0,t.jsx)(s.p,{children:"When connecting to remote servers, ensure the usage of DNS instead of IP addresses to avoid issues."}),"\n",(0,t.jsx)(s.p,{children:"The following attributes are evaluated:"}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-yaml",children:"- http.url\n- db.connection_string\n"})}),"\n",(0,t.jsxs)(s.p,{children:["If span kind is ",(0,t.jsx)(s.code,{children:'"client"'}),", the following attributes are evaluated:"]}),"\n",(0,t.jsx)(s.pre,{children:(0,t.jsx)(s.code,{className:"language-yaml",children:"- net.peer.name\n"})}),"\n",(0,t.jsx)(s.h2,{id:"options",children:"Options"}),"\n",(0,t.jsx)(s.p,{children:"This rule has the following options:"}),"\n",(0,t.jsxs)(s.ul,{children:["\n",(0,t.jsxs)(s.li,{children:[(0,t.jsx)(s.code,{children:'"error"'})," requires DNS over IP addresses"]}),"\n",(0,t.jsxs)(s.li,{children:[(0,t.jsx)(s.code,{children:'"disabled"'})," disables the DNS over IP addresses verification"]}),"\n",(0,t.jsxs)(s.li,{children:[(0,t.jsx)(s.code,{children:'"warning"'})," verifies DNS over IP addresses but does not impact the analyzer score"]}),"\n"]}),"\n",(0,t.jsx)(s.h2,{id:"when-not-to-use-it",children:"When Not To Use It"}),"\n",(0,t.jsx)(s.p,{children:"If you intentionally use and record IP addresses then you can disable this rule."})]})}function u(e={}){const{wrapper:s}={...(0,r.R)(),...e.components};return s?(0,t.jsx)(s,{...e,children:(0,t.jsx)(c,{...e})}):c(e)}},28453:(e,s,n)=>{n.d(s,{R:()=>o,x:()=>a});var t=n(96540);const r={},i=t.createContext(r);function o(e){const s=t.useContext(i);return t.useMemo((function(){return"function"==typeof e?e(s):{...s,...e}}),[s,e])}function a(e){let s;return s=e.disableParentContext?"function"==typeof e.components?e.components(r):e.components||r:o(e.components),t.createElement(i.Provider,{value:s},e.children)}}}]);