"use strict";(self.webpackChunktracetest_docs=self.webpackChunktracetest_docs||[]).push([[2155],{24549:(e,t,n)=>{n.r(t),n.d(t,{assets:()=>c,contentTitle:()=>i,default:()=>h,frontMatter:()=>r,metadata:()=>o,toc:()=>l});var s=n(74848),a=n(28453);const r={id:"expressions",title:"Expressions",description:"Expressions are used to add values that are only known during execution time. Build integration and end-to-end tests with OpenTelemetry traces with Tracetest.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},i=void 0,o={id:"concepts/expressions",title:"Expressions",description:"Expressions are used to add values that are only known during execution time. Build integration and end-to-end tests with OpenTelemetry traces with Tracetest.",source:"@site/docs/concepts/expressions.mdx",sourceDirName:"concepts",slug:"/concepts/expressions",permalink:"/concepts/expressions",draft:!1,unlisted:!1,editUrl:"https://github.com/kubeshop/tracetest/blob/main/docs/docs/concepts/expressions.mdx",tags:[],version:"current",frontMatter:{id:"expressions",title:"Expressions",description:"Expressions are used to add values that are only known during execution time. Build integration and end-to-end tests with OpenTelemetry traces with Tracetest.",hide_table_of_contents:!1,keywords:["tracetest","trace-based testing","observability","distributed tracing","testing"],image:"https://res.cloudinary.com/djwdcmwdz/image/upload/v1698686403/docs/Blog_Thumbnail_14_rsvkmo.jpg"},sidebar:"tutorialSidebar",previous:{title:"Selectors",permalink:"/concepts/selectors"},next:{title:"Assertions",permalink:"/concepts/assertions"}},c={},l=[{value:"Features",id:"features",level:2},{value:"Reference Span Attributes",id:"reference-span-attributes",level:3},{value:"Reference Variables",id:"reference-variables",level:3},{value:"Arithmetic Operations",id:"arithmetic-operations",level:3},{value:"String Interpolation",id:"string-interpolation",level:3},{value:"Filters",id:"filters",level:3},{value:"JSON Path",id:"json-path",level:4},{value:"RegEx",id:"regex",level:4},{value:"RegEx Group",id:"regex-group",level:4},{value:"Get Index",id:"get-index",level:3},{value:"Length",id:"length",level:3},{value:"Type",id:"type",level:3},{value:"JSON Comparison",id:"json-comparison",level:3}];function d(e){const t={a:"a",code:"code",h2:"h2",h3:"h3",h4:"h4",li:"li",p:"p",pre:"pre",ul:"ul",...(0,a.R)(),...e.components};return(0,s.jsxs)(s.Fragment,{children:[(0,s.jsx)(t.p,{children:"Tracetest allows you to add expressions when writing your tests. They are a nice and clean way of adding values that are only known during execution time. For example, when referencing a variable, a span attribute or even arithmetic operations."}),"\n",(0,s.jsx)(t.h2,{id:"features",children:"Features"}),"\n",(0,s.jsxs)(t.ul,{children:["\n",(0,s.jsx)(t.li,{children:"Reference span attributes"}),"\n",(0,s.jsx)(t.li,{children:"Reference variables"}),"\n",(0,s.jsx)(t.li,{children:"Arithmetic operations"}),"\n",(0,s.jsx)(t.li,{children:"String interpolation"}),"\n",(0,s.jsx)(t.li,{children:"Filters"}),"\n",(0,s.jsx)(t.li,{children:"JSON comparison"}),"\n"]}),"\n",(0,s.jsx)(t.h3,{id:"reference-span-attributes",children:"Reference Span Attributes"}),"\n",(0,s.jsxs)(t.p,{children:["When building assertions, you might need to assert if a certain span contains an attribute and that this attribute has a specific value. To accomplish this with Tracetest, you can use expressions to get the value of the span. When referencing an attribute, add the prefix ",(0,s.jsx)(t.code,{children:"attr:"})," and its name. For example, imagine you have to check if the attribute ",(0,s.jsx)(t.code,{children:"service.name"})," is equal to ",(0,s.jsx)(t.code,{children:"cart-api"}),". Use the following statement:"]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'attr:service.name = "cart-api"\n'})}),"\n",(0,s.jsx)(t.h3,{id:"reference-variables",children:"Reference Variables"}),"\n",(0,s.jsxs)(t.p,{children:["Create variables in Tracetest based on the trace obtained by the test to enable assertions that require values from other spans. Variables use the prefix ",(0,s.jsx)(t.code,{children:"var:"})," and its name. For example, a variable called ",(0,s.jsx)(t.code,{children:"user_id"})," would be referenced as ",(0,s.jsx)(t.code,{children:"var:user_id"})," in an expression."]}),"\n",(0,s.jsx)(t.h3,{id:"arithmetic-operations",children:"Arithmetic Operations"}),"\n",(0,s.jsx)(t.p,{children:"Sometimes we need to manipulate data to ensure our test data is correct. As an example, we will use a purchase operation. How you would make sure that, after the purchase, the product inventory is smaller than before? For this, we might want to use arithmetic operations:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:"attr:product.stock = attr:product.stok_before_purchase - attr:product.number_bought_items\n"})}),"\n",(0,s.jsx)(t.h3,{id:"string-interpolation",children:"String Interpolation"}),"\n",(0,s.jsx)(t.p,{children:"Some tests might require strings to be compared, but maybe you need to generate a dynamic string that relies on a dynamic value. This might be used in an assertion or even in the request body referencing a variable."}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'attr:error.message = "Could not withdraw ${attr:withdraw.amount}, your balance is insufficient."\n'})}),"\n",(0,s.jsxs)(t.p,{children:["Note that within ",(0,s.jsx)(t.code,{children:"${}"})," you can add any expression, including arithmetic operations and filters."]}),"\n",(0,s.jsx)(t.h3,{id:"filters",children:"Filters"}),"\n",(0,s.jsx)(t.p,{children:"Filters are functions that are executed using the value obtained by the expression. They are useful to transform the data. Multiple filters can be chained together. The output of the previous filter will be used as the input to the next until all filters are executed."}),"\n",(0,s.jsx)(t.h4,{id:"json-path",children:"JSON Path"}),"\n",(0,s.jsx)(t.p,{children:"This filter allows you to filter a JSON string and obtain only data that is relevant."}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'\'{ "name": "Jorge", "age": 27, "email": "jorge@company.com" }\' | json_path \'.age\' = 27\n'})}),"\n",(0,s.jsx)(t.p,{children:"If multiple values are matched, the output will be a flat array containing all values."}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'\'{ "array": [{"name": "Jorge", "age": 27}, {"name": "Tim", "age": 52}]}\' | json_path \'$.array[*]..["name", "age"] = \'["Jorge", 27, "Tim", 52]\'\n'})}),"\n",(0,s.jsx)(t.h4,{id:"regex",children:"RegEx"}),"\n",(0,s.jsx)(t.p,{children:"Filters part of the input that match a RegEx. Imagine you have a specific part of a text that you want to extract:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:"'My account balance is $48.52' | regex '\\$\\d+(\\.\\d+)?' = '$48.52'\n"})}),"\n",(0,s.jsx)(t.h4,{id:"regex-group",children:"RegEx Group"}),"\n",(0,s.jsx)(t.p,{children:"If matching more than one value is required, you can define groups for your RegEx and extract multiple values at once."}),"\n",(0,s.jsx)(t.p,{children:"Wrap the groups you want to extract with parentheses."}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:"'Hello Marcus, today you have 8 meetings' | regex_group 'Hello (\\w+), today you have (\\d+) meetings' = '[\"Marcus\", 8]'\n"})}),"\n",(0,s.jsx)(t.h3,{id:"get-index",children:"Get Index"}),"\n",(0,s.jsx)(t.p,{children:"Some filters might result in an array. If you want to assert just part of this array, this filter allows you to pick one element from the array based on its index."}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:"'{ \"array\": [1, 2, 3] }' | json_path '$.array[*]' | get_index 1 = 2\n"})}),"\n",(0,s.jsxs)(t.p,{children:["You can select the last item from a list by specifying ",(0,s.jsx)(t.code,{children:"'last'"})," as the argument for the filter:"]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:"'{ \"array\": [1, 2, 3] }' | json_path '$.array[*]' | get_index 'last' = 3\n"})}),"\n",(0,s.jsx)(t.h3,{id:"length",children:"Length"}),"\n",(0,s.jsxs)(t.p,{children:["Returns the size of the input array. If it's a single value, it will return 1. Otherwise it will return ",(0,s.jsx)(t.code,{children:"length(input_array)"}),"."]}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:"'{ \"array\": [1, 2, 3] }' | json_path '$.array[*]' | length = 3\n"})}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'"my string" | length = 1\n'})}),"\n",(0,s.jsx)(t.h3,{id:"type",children:"Type"}),"\n",(0,s.jsx)(t.p,{children:"Return the type of the input as a string."}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'2 | type = "number"\n'})}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'2s | type = "duration"\n'})}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'"something" | type = "string"\n'})}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'[1, 2s, "this is a string"] | type = "array"\n'})}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'attr:myapp.operations | get_index 2 | type = "string"\n'})}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'# check if attribute is either a number of a string\n["number", "string"] contains attr:my_attribute | type\n'})}),"\n",(0,s.jsx)(t.h3,{id:"json-comparison",children:"JSON Comparison"}),"\n",(0,s.jsxs)(t.p,{children:["When working with APIs, it's very common for developers to check if the response contains data while ignoring all the noise a response can have. For this, you can use the ",(0,s.jsx)(t.code,{children:"contains"})," comparator in JSON objects. It works just like ",(0,s.jsx)(t.a,{href:"https://jestjs.io/pt-BR/docs/expect#tomatchobjectobject",children:"Jest's toMatchObject"}),"."]}),"\n",(0,s.jsx)(t.p,{children:"The order of the attributes doesn't matter and the left side of the expression can contain more attributes than the right side."}),"\n",(0,s.jsx)(t.p,{children:"Some examples:"}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'# If user is 32 years old, it passes\n\'{"name": "john", "age": 32}\' contains \'{"age": 32}\'\n'})}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'# If any of the users is called "maria", it passes\n\'[{"name": "john", "age": 32}, {"name": "maria", "age": 63}]\' contains \'{"name": "maria"}\'\n'})}),"\n",(0,s.jsx)(t.pre,{children:(0,s.jsx)(t.code,{className:"language-css",children:'# In this case, both ages must be part of the JSON array\n\'[{"name": "john", "age": 32}, {"name": "maria", "age": 63}]\' contains \'[{"age": 63}, {"age": 32}]\'\n'})})]})}function h(e={}){const{wrapper:t}={...(0,a.R)(),...e.components};return t?(0,s.jsx)(t,{...e,children:(0,s.jsx)(d,{...e})}):d(e)}},28453:(e,t,n)=>{n.d(t,{R:()=>i,x:()=>o});var s=n(96540);const a={},r=s.createContext(a);function i(e){const t=s.useContext(r);return s.useMemo((function(){return"function"==typeof e?e(t):{...t,...e}}),[t,e])}function o(e){let t;return t=e.disableParentContext?"function"==typeof e.components?e.components(a):e.components||a:i(e.components),s.createElement(r.Provider,{value:t},e.children)}}}]);