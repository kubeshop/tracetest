<p align="center">  
  <img style="width:66%" src="assets/tracetest-logo-color-w-white-text.svg#gh-dark-mode-only" alt="Tracetest Logo Light"/>
  <img style="width:66%" src="assets/tracetest-logo-color-w-black-text.svg#gh-light-mode-only" alt="Tracetest Logo Dark" />
</p>

<p align="center">
  End-to-end tests powered by OpenTelemetry. For QA, Dev, & Ops.
</p>

<p align="center">
  <!--<a href="https://tracetest.io">Website</a>&nbsp;|&nbsp; -->
  <a href="https://github.com/kubeshop/tracetest#try-the-demo--give-us-feedback">Live Demo</a>&nbsp;|&nbsp; 
  <a href="https://kubeshop.github.io/tracetest">Documentation</a>&nbsp;|&nbsp; 
  <a href="https://twitter.com/tracetest_io">Twitter</a>&nbsp;|&nbsp; 
  <a href="https://discord.gg/eBvEQRVyKX">Discord</a>&nbsp;|&nbsp; 
  <a href="https://kubeshop.io/category/tracetest">Blog (TBD!)</a>
</p>

<p align="center">
  <a href="https://github.com/kubeshop/tracetest/releases"><img title="Release" src="https://img.shields.io/github/v/release/kubeshop/tracetest"/></a>
  <a href=""><img title="Docker builds" src="https://img.shields.io/docker/automated/kubeshop/tracetest"/></a>
  <a href="https://github.com/kubeshop/tracetest/releases"><img title="Release date" src="https://img.shields.io/github/release-date/kubeshop/tracetest"/></a>
</p>

<p align="center">
  <a target="_new" href="https://www.youtube.com/watch?v=WMRicNlaehc">
    <img src="https://img.youtube.com/vi/WMRicNlaehc/0.jpg" style="width:66%;height:auto">
    <p align="center">
      Click on the image or this link to watch the "Tracetest Intro Video" video (7 mins)
    </p>
  </a>
</p>

# Overview

Testing and debugging software built on microservices architectures is not an easy task. Multiple services, teams, programming languages, and technologies are involved. We want to help you write tests across all this complexity by leveraging your existing investment in OpenTelemetry tracing.

Tracetest makes it easy:

1. Pick an API to test.
2. Run a test, and get the trace.
3. The trace is the blueprint of your system under test. It shows all the steps the system has taken to execute the request.
4. Use this blueprint to define assertions through Tracetest UI.
5. Add assertions on different services, checking return statuses, data, or even execution times of a system.
6. Run the tests.

![Assertions](/assets/assertions.png)

Once the test is built, it can be run automatically as part of a build process. Every test has a trace attached, allowing you to immediately see what worked, and what did not, reducing the need to reproduce the problem to see the underlying issue.

# System Diagram

<div style="text-align:center;"><img src="/assets/tracetest-diagram-01.png"></div>

# Try the demo & give us feedback

We have a live demo environment with a couple systems you can test against. Use the 'Try The Demo' button below to launch it. You will need to know the urls you can test against - here are some examples that work:

| System               | Description      | URL                                                                   | Method | Request Body                                                                                                                           |
| -------------------- | ---------------- | --------------------------------------------------------------------- | ------ | -------------------------------------------------------------------------------------------------------------------------------------- |
| Shopping app         | Generic get      | http://shop/buy                                                       | GET    |
| Pokemon Microservice | Get a Pokemon    | http://demo-pokemon-api.demo.svc.cluster.local/pokemon?take=20&skip=0 | GET    |
| Pokemon Microservice | Add a Pokemon    | http://demo-pokemon-api.demo.svc.cluster.local/pokemon                | POST   | { "name": "meowth", "type": "normal","imageUrl": "https://assets.pokemon.com/assets/cms2/img/pokedex/full/052.png","isFeatured": true} |
| Pokemon Microservice | Import a Pokemon | http://demo-pokemon-api.demo.svc.cluster.local/pokemon/import         | POST   | { "id": 52 }                                                                                                                           |

More documentation about the installed Pokemon Microservice App, PMA, is at [Pokemon Microservice App github](https://github.com/kubeshop/pokeshop/blob/master/docs/overview.md)

Wanna play with it? [Try the Live Demo](https://demo.tracetest.io/)

[![button](/assets/button_try_tracetest.png)](https://demo.tracetest.io/)

We’re looking for feedback to help make Tracetest even better for developers, QA testers, and DevOPs. Please give us feedback on [Discord](https://discord.gg/eBvEQRVyKX) or [create an issue on Github](https://github.com/kubeshop/tracetest/issues/new/choose)

# Getting Started

Check out the [Installation](https://kubeshop.github.io/tracetest/installing/) and
[Getting Started](https://kubeshop.github.io/tracetest/getting-started/) guides to set up Tracetest and
run your first tests! It is still a 'work in progress' so please provide us with any and all [feedback](https://github.com/kubeshop/tracetest/issues/new/choose) - we live for input and will respond!

We are in the our early days with the project and need your help. Have an idea to improve it? Please Create an issue here or join our community on Discord (link).

Follow us on [Twitter at @tracetest_io](https://twitter.com/tracetest_io) for updates

Give us a star on Github if you're interested in the project!

# Documentation

Is available at [https://kubeshop.github.io/tracetest](https://kubeshop.github.io/tracetest)
