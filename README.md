# Tracetest

## Overview

Testing and debugging software built on Micro-Services architectures is hard.

As many as 30 to 100 services may be involved in a single flow. Written in multiple languages. With several backend data stores, message busses, and technologies. Understanding the flow is hard - having enough experience & wide ranging knowledge to create tests to verify it is working properly is even harder.

Tracetest makes this easy. Pick an api to test. Tracetest uses your tracing infrastructure to trace this api call. This trace is the blueprint of your entire system, showing all the activity. Use this blueprint to graphically define assertions on different services throughout the trace, checking return statuses, data, or even execution times of systems.

Examples:
assert that all database calls return in less than 250 ms
assert that one particular micro service returns a 200 code when called
assert that a Kafka queue successful delivers a payload to a dependent micro service.

Once the test is built, it can be run automatically as part of a build process or manually. Every test has a trace attached, allowing you to immediately see what worked, and what did not, reducing the need to reproduce the problem to see the underlying issue.

# Development

## Web Folder Structure

### Web Folder Tree
 
* [public/](./web/public)
* [src/](./web/src)
  * [components/](./web/src/components)
  * [hooks/](./web/src/hooks)
  * [lib/](./web/src/lib)
  * [navigation/](./web/src/navigation)
  * [pages/](./web/src/pages)
  * [redux/](./web/src/redux)
  * [services/](./web/src/services)
  * [types/](./web/src/types)
  * [utils/](./web/src/utils)

| Folder  |  Description  |
|---|---|
| public/ | contains html and can put any scripts, or static files  |
|  components/ |  any reusable components  |
| hooks/  | any reusable hooks  |
| lib/  |  any constants  |
| navigation/  |  react-router setup |
| pages/  |  each page has its folder that contains styled file, and sub components |
| redux/  | state management setup  |
| services/  |  any class or functions that talk with data and processing it |
| types/| defined models used in the web app| |
|utils/| any pure functions that shared in the app ||
