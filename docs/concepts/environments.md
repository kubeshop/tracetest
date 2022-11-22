# Environments

A common use case for tests is to assert the same behavior across multiple environments (dev, staging, and production, for example). But to make sure all environments will have the same behavior, it's important that the tests executed against those environments test the same aspects. To reduce risks of diverging tests, Tracetest allows you to organize different environments configurations using global objects called Environments.

## How environments work

Environments are objects containing variables that can be referenced by tests. You can have the same test and just swap the information for each environment by changing which environment object to be used to run the test. To illustrate it, consider an app that is deployed in three stages: `dev`, `staging`, and `production`. We can execute the same test in all those environments, however, both `URL` and `credentials` change from environment to environment. To be able to run the same test against the three deployments of the app, you can create three environments:

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

Both `env:URL` and `env:API_TOKEN` would be replaced by the variables defined in the chosen environment to be used when running the test. So, if the chosen environment was `dev.env`, its values would be replaced by `https://app-dev.com` and `dev-key` respectively.

