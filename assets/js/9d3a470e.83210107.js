"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[1102],{7239:(e,t,i)=>{i.r(t),i.d(t,{assets:()=>c,contentTitle:()=>r,default:()=>u,frontMatter:()=>s,metadata:()=>o,toc:()=>l});var n=i(74848),a=i(28453);const s={id:"creating-variable-sets",title:"Defining Variable Sets as Text Files",description:"Configure Variable Sets in Tracetest by using a configuration file that can be applied to your Tracetest instance.",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},r=void 0,o={id:"cli/creating-variable-sets",title:"Defining Variable Sets as Text Files",description:"Configure Variable Sets in Tracetest by using a configuration file that can be applied to your Tracetest instance.",source:"@site/docs/cli/creating-variable-sets.mdx",sourceDirName:"cli",slug:"/cli/creating-variable-sets",permalink:"/cli/creating-variable-sets",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/cli/creating-variable-sets.mdx",tags:[],version:"current",frontMatter:{id:"creating-variable-sets",title:"Defining Variable Sets as Text Files",description:"Configure Variable Sets in Tracetest by using a configuration file that can be applied to your Tracetest instance.",keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"tutorialSidebar",previous:{title:"Defining Monitors as Text Files",permalink:"/cli/configuring-monitors"},next:{title:"CLI Reference",permalink:"/cli/reference/tracetest"}},c={},l=[{value:"Using Secrets in Variable Sets",id:"using-secrets-in-variable-sets",level:2}];function d(e){const t={a:"a",admonition:"admonition",blockquote:"blockquote",code:"code",h2:"h2",mdxAdmonitionTitle:"mdxAdmonitionTitle",p:"p",pre:"pre",...(0,a.R)(),...e.components};return(0,n.jsxs)(n.Fragment,{children:[(0,n.jsx)(t.p,{children:"This page showcases how to create and edit variable sets with the CLI."}),"\n",(0,n.jsxs)(t.admonition,{type:"info",children:[(0,n.jsx)(t.mdxAdmonitionTitle,{}),(0,n.jsx)(t.p,{children:"For details on creating and editing variable sets in the Web UI, please visit Web UI Creating Variable Sets."})]}),"\n",(0,n.jsx)(t.admonition,{type:"tip",children:(0,n.jsx)(t.p,{children:(0,n.jsx)(t.a,{href:"/concepts/variable-sets",children:"To read more about variable sets check out variable sets concepts."})})}),"\n",(0,n.jsx)(t.p,{children:"Just like Data Stores, you can also manage your variable sets using the CLI and definition files."}),"\n",(0,n.jsx)(t.p,{children:"A definition file looks like the following:"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-yaml",children:"type: VariableSet\nspec:\n    name: Production\n    description: Production env variables\n    values:\n    - key: URL\n      value: https://app-production.company.com\n    - key: API_KEY\n      value: mysecret\n"})}),"\n",(0,n.jsxs)(t.p,{children:["In order to apply this configuration to your Tracetest instance, make sure to have your ",(0,n.jsx)(t.a,{href:"/cli/configuring-your-cli",children:"CLI configured"})," and run:"]}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{className:"language-sh",children:"tracetest apply variableset -f <variableset.yaml>\n"})}),"\n",(0,n.jsxs)(t.blockquote,{children:["\n",(0,n.jsxs)(t.p,{children:["If the file contains the property ",(0,n.jsx)(t.code,{children:"spec.id"}),", the operation will be considered a variable set update. If you try to apply a variable set and you get the error: ",(0,n.jsx)(t.code,{children:"could not apply variableset: 404 Not Found"}),", it means the provided ID doesn't exist. Either update the ID to reference an existing variable set or remove the property from the file. Tracetest will create a new variable set and a new ID."]}),"\n"]}),"\n",(0,n.jsx)(t.h2,{id:"using-secrets-in-variable-sets",children:"Using Secrets in Variable Sets"}),"\n",(0,n.jsx)(t.p,{children:"A VariableSet with a secret values can be registered via the CLI like this:"}),"\n",(0,n.jsx)(t.pre,{children:(0,n.jsx)(t.code,{children:"type: VariableSet\nspec:\n  id: test-1\n  name: Test 1\n  description: test\n  values:\n  - key: SOME_SECRET\n    value: my-precious\n    type: secret\n  - key: NO_SO_SECRET\n    value: aha!\n    type: raw\n"})})]})}function u(e={}){const{wrapper:t}={...(0,a.R)(),...e.components};return t?(0,n.jsx)(t,{...e,children:(0,n.jsx)(d,{...e})}):d(e)}},28453:(e,t,i)=>{i.d(t,{R:()=>r,x:()=>o});var n=i(96540);const a={},s=n.createContext(a);function r(e){const t=n.useContext(s);return n.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function o(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(a):e.components||a:r(e.components),n.createElement(s.Provider,{value:t},e.children)}}}]);