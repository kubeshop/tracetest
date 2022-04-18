<p align="center">  
  <img style="width:66%" src="assets/tracetest-color-white.png#gh-dark-mode-only" alt="Tracetest Logo Light"/>
  <img style="width:66%" src="assets/tracetest-color-dark.png#gh-light-mode-only" alt="Tracetest Logo Dark" />
</p>

<p align="center">
  End-to-end tests powered by OpenTelemetry. For QA, Dev, & Ops.
</p>

<p align="center">
  <!--<a href="https://tracetest.io">Website</a>&nbsp;|&nbsp; -->
  <a href="https://demo.tracetest.io">Live Demo</a>&nbsp;|&nbsp; 
  <a href="https://kubeshop.github.io/tracetest">Documentation</a>&nbsp;|&nbsp; 
  <a href="https://twitter.com/tracetest_io">Twitter</a>&nbsp;|&nbsp; 
  <a href="https://discord.gg/eBvEQRVyKX">Discord</a>&nbsp;|&nbsp; 
  <a href="https://kubeshop.io/category/tracetest">Blog (TBD!)</a>
</p>

<p align="center">
  <a href="https://github.com/kubeshop/tracetest/releases"><img title="Release" src="https://img.shields.io/github/v/release/kubeshop/tracetest"/></a>
  <a href=""><img title="Downloads" src="https://img.shields.io/github/downloads/kubeshop/tracetest/total.svg"/></a>
  <a href=""><img title="Go version" src="https://img.shields.io/github/go-mod/go-version/kubeshop/tracetest"/></a>
  <a href=""><img title="Docker builds" src="https://img.shields.io/docker/automated/kubeshop/tracetest"/></a>
  <a href="https://github.com/kubeshop/tracetest/releases"><img title="Release date" src="https://img.shields.io/github/release-date/kubeshop/tracetest"/></a>
</p>

<p align="center">
  <a target="_new" href="https://www.youtube.com/watch?v=GVvgLuxdrXE&t=47s">
    <img src="assets/intro-to-tracetest.jpg" style="width:66%;height:auto">
    <p align="center">
      Click on the image or this link to watch the "Intro to Tracetest" short video (3 mins)
    </p>
  </a>
</p>

# Overview

Testing and debugging software built on Micro-Services architectures is not an easy task. Multiple services, different teams, various programming languages and  technologies involved. We would like to help you write tests across all this complexity.

Key Value Prop = Tracetest uses Open Telemetery tracing infrastructure to ....
Tracetest makes it easy. For example, pick an API to test. Get its trace. This trace is the blueprint of your system (or of that API?}, showing all the steps. Use this blueprint to graphically define assertions through Tracetest UI on different services throughout the trace, checking return statuses, data, or even execution times of a system (system or API).

![Assertions](/assets/assertions.png)

Once the test is built, it can be run automatically as part of a build process. Every test has a trace attached, allowing you to immediately see what worked, and what did not, reducing the need to reproduce the problem to see the underlying issue.

# System Diagram

<div style="text-align:center;"><img src="/assets/tracetest-diagram-01.png"></div>

# Try the demo & give us feedback

Wanna play with it? <button name="button" onClick="https://demo.tracetest.io">Try out the live Tracetest demo!</button>

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
