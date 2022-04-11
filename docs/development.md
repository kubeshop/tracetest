# Development

## Web Folder Structure

### Web Folder Tree

- [public/](./web/public)
- [src/](./web/src)
  - [components/](./web/src/components)
  - [hooks/](./web/src/hooks)
  - [lib/](./web/src/lib)
  - [navigation/](./web/src/navigation)
  - [pages/](./web/src/pages)
  - [redux/](./web/src/redux)
  - [services/](./web/src/services)
  - [types/](./web/src/types)
  - [utils/](./web/src/utils)

| Folder      | Description                                                            |
| ----------- | ---------------------------------------------------------------------- | --- |
| public/     | contains html and can put any scripts, or static files                 |
| components/ | any reusable components                                                |
| hooks/      | any reusable hooks                                                     |
| lib/        | any constants                                                          |
| navigation/ | react-router setup                                                     |
| pages/      | each page has its folder that contains styled file, and sub components |
| redux/      | state management setup                                                 |
| services/   | any class or functions that talk with data and processing it           |
| types/      | defined models used in the web app                                     |     |
| utils/      | any pure functions that shared in the app                              |     |
