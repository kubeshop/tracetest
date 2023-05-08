# Tracetest internal testing

In this folder, we have some of the test automation structures used to evaluate Tracetest. 

The tests are:
- **cli-e2etest**: CLI End-to-end tests, where we run a compiled CLI against some setups of Tracetest infrastructure and we check if the CLI commands are properly working;
- **cli-smoketest**: simple CLI test, where we check the CLI was correctly compiled and can run simple commands (as `tracetest version`);
- **server-tracetesting**: set of [dogfooding](https://en.wikipedia.org/wiki/Eating_your_own_dog_food) tests, where run some trace-based tests against the current version of Tracetest to check if the Tracetest API is working fine.
