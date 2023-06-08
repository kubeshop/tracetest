# Tracetest Analyzer Concepts

The Tracetest Analyzer is a plugin-based framework used to analyze OpenTelemetry traces to help teams improve their instrumentation data, find potential problems and provide tips to fix the problems.

## Problem

Today, implementing open telemetry instrumentation in any application can become overwhelming as most of the documentation is not centralized, can be confusing and as the problem is relatively new there aren’t many good tutorials or guides to follow.

Usually, there are official libraries for the most common languages that can be used to add basic auto-instrumentation. Still, when it comes to adding custom instrumentation it can be confusing to understand if what is added is aligned to the Otel standard in terms of having the right information and/or not leaking any sensitive data.

Another problem is that adding instrumentation is the very first step to achieving and improving an application architecture, having visibility is just a way to understand what’s going on, but then it can become a little more difficult to understand what's next. 

## Solution

Having a linting rule that is fully plugin-based to evaluate Tracetest to find problems. Allowing simple ways for users to understand the quality of their instrumentation data across the whole framework, catching possible security breaches, and providing tips and guidance on how to fix them.

### Concepts

**Plugin**

Is the encapsulation of an `N` number of rules, with a name and a category that defines its specific goal.

**Rule**

Contains the set of validations to be evaluated using one or multiple traces depending on the rule type. The result when applying a rule should be if it passed or failed, displaying extra information using tips for user guidance.

Each rule has a name, description, and logic to evaluate the whole trace or a specific span.

**Rule Types**

There are two main rule types:

- Multi trace. Requires a historic number of traces to identify and evaluate the rule.
- Single Trace. Encapsulated to the current trace, no external data is required.

NOTE: This documentation will be focused on single trace rules for timing purposes.

**Linter Resource**

Allows the user to configure the linter framework, a global on/off switch, an opt-in/out for specific plugins or rules, setting up required and optional rules, etc.


