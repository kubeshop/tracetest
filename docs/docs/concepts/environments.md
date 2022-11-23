# Environments

A common use case for tests is to assert the same behavior across multiple environments (dev, staging, and production, for example). To make sure all environments will have the same behavior, it is important that the tests executed against those environments test the same aspects. To reduce the risks of diverging tests, Tracetest allows you to organize different environments configurations using global objects called **Environments**.

## How Environments Work

Environments are objects containing variables that can be referenced by tests. You can use a single test and provide the information on which environment object will be used to run the test. To illustrate this, consider an app that is deployed in three stages: `dev`, `staging`, and `production`. We can execute the same test in all those environments, however, both `URL` and `credentials` change from environment to environment. To run the same test against the three deployments of the app, you can create three environments:

```dotenv
# dev.env
URL=https://app-dev.com
API_TOKEN=dev-key
```

```dotenv
# staging.env
URL=https://app-staging.com
API_TOKEN=staging-key
```

```dotenv
# production.env
URL=https://app-prod.com
API_TOKEN=prod-key
```

Now consider the following test:

```yaml
type: Test
specs:
  name: Test user creation
  trigger:
    type: http
    httpRequest:
        url: "${env:URL}/api/users"
        method: POST
        body: '{}'
        authentication:
          type: bearer
          bearer:
            token: "${env:API_TOKEN}"
```

Both `env:URL` and `env:API_TOKEN` would be replaced by the variables defined in the chosen environment where the test will run. So, if the chosen environment was `dev.env`, its values would be replaced by `https://app-dev.com` and `dev-key`, respectively.
