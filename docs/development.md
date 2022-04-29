# Development Guide

The phiolosphy used to build the frontend is based on object composition and functional programming principles, as well as following the idea of "code that is updated together should be together" as much as we can. To achieve this the main funcionality is divided into the following concepts:

**Gateways**
Interfaces to the outbound backend services and the browser.
Used by  Components, Services.

**Models**
Encapsulate a set of attributes and modules.
Used by Services, APIs.

**Services**
Define and drive the Business logic functionality
Used by Components, Selectors.

**Factories**
Component Facades that allow rendering different versions depending on the use case.
Used by Components.

**Selectors**
Parse, aggregate and clean data across the app.
Used by Components.

### Styling Framework
For styling and web components the main framework in use is [Ant Design](https://ant.design/).

### Web Folder Structure

####  Web Folder Tree

- [public/](./web/public)
- [src/](./web/src)
  - [components/](./web/src/components)
  - [gateways/](./web/src/gateways)
  - [assets/](./web/src/assets)
  - [hooks/](./web/src/hooks)
  - [constants/](./web/src/lib)
  - [entities/](./web/src/entities/)
  - [navigation/](./web/src/navigation)
  - [pages/](./web/src/pages)
  - [redux/](./web/src/redux)
  - [types/](./web/src/types)
  - [utils/](./web/src/utils)

| Folder      | Description                                                            |
| ----------- | ---------------------------------------------------------------------- |
| public/     | contains html and can put any scripts, or static files                 |
| assets/     | contains icons, images or any other private non-code file used in the app                 |
| gateways/ | any API connector                                                |
| components/ | any reusable components                                                |
| hooks/      | any reusable hooks                                                     |
| constants/  | any constants                                                          |
| navigation/ | react-router setup                                                     |
| pages/      | each page has its folder that contains styled file, and sub components |
| redux/      | state management setup                                                 |
| entities/   | each entity from the application has its own folder containing services, models, gateways, etc           |
| types/      | defined models used in the web app                                     |
| utils/      | any pure functions that shared in the app                              |

#### File Types

| Type      | Exstension                | Example
| ----------- | ----------------------- | ----
| Service     | .service.ts             | Span.service.ts
| Model      | .model.ts               | Span.model.ts
| Factory   | .factory.tsx             | Diagram.factory.tsx
| Component | .tsx                     | Input.tsx
| Styled Components | .styled.ts(x)    | HomePage.styled.ts
| Selectors | .selectors.ts | Span.selectors.ts
| Gateways | .gateway.ts    | TestApi.gateway.ts
| Test | .test.ts | Span.test.ts|
| Constants | .constants.ts | Trace.constants.ts
| Types | .types.ts | Assertion.types.ts
