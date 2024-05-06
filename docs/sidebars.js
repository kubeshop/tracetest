/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

// generated docs pages
/** @type {import('@docusaurus/plugin-content-docs/lib/sidebars/types').SidebarItem[]} */
const cliGeneratedPages = require("./docs/cli/reference/cli-sidebar");

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  // tutorialSidebar: [{type: 'autogenerated', dirName: '.'}],

  // But you can create a sidebar manually

  liveExamplesSidebar: [
    {
      type: "category",
      label: "OpenTelemetry Store Demo",
      link: {
        type: "doc",
        id: "live-examples/opentelemetry-store/overview",
      },
      items: [
        {
          type: "category",
          label: "Use Cases",
          items: [
            {
              type: "doc",
              id: "live-examples/opentelemetry-store/use-cases/add-item-into-shopping-cart",
              label: "Add Item into Shopping Cart",
            },
            {
              type: "doc",
              id: "live-examples/opentelemetry-store/use-cases/check-shopping-cart-contents",
              label: "Check Shopping Cart Contents",
            },
            {
              type: "doc",
              id: "live-examples/opentelemetry-store/use-cases/checkout",
              label: "Checkout",
            },
            {
              type: "doc",
              id: "live-examples/opentelemetry-store/use-cases/get-recommended-products",
              label: "Get Recommended Products",
            },
            {
              type: "doc",
              id: "live-examples/opentelemetry-store/use-cases/user-purchasing-products",
              label: "User Purchasing Products",
            },
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Pokeshop API Demo",
      link: {
        type: "doc",
        id: "live-examples/pokeshop/overview",
      },
      items: [
        {
          type: "category",
          label: "Use Cases",
          items: [
            {
              type: "doc",
              id: "live-examples/pokeshop/use-cases/add-pokemon",
              label: "Add Pokemon",
            },
            {
              type: "doc",
              id: "live-examples/pokeshop/use-cases/list-pokemon",
              label: "List Pokemon",
            },
            {
              type: "doc",
              id: "live-examples/pokeshop/use-cases/get-pokemon-by-id",
              label: "Get Pokemon by ID",
            },
            {
              type: "doc",
              id: "live-examples/pokeshop/use-cases/import-pokemon",
              label: "Import Pokemon",
            },
            {
              type: "doc",
              id: "live-examples/pokeshop/use-cases/import-pokemon-from-stream",
              label: "Import Pokemon from Stream",
            },
          ],
        },
      ],
    },
  ],

  examplesTutorialsSidebar: [
    {
      type: "doc",
      id: "examples-tutorials/overview",
      label: "Overview",
    },
    {
      type: "category",
      link: {
        type: "doc",
        id: "examples-tutorials/recipes",
      },
      label: "Recipes",
      items: [
        {
          type: "category",
          label: "Streams and Message Queues",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/testing-kafka-go-api-with-opentelemetry-tracetest",
              label: "Testing Kafka",
            },
          ],
        },
        {
          type: "category",
          label: "Serverless Functions",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/testing-vercel-functions-with-opentelemetry-tracetest",
              label: "Testing Vercel Functions",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/testing-lambda-functions-with-opentelemetry-tracetest",
              label: "Testing AWS Lambda Functions",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/testing-cloudflare-workers-with-opentelemetry-tracetest",
              label: "Testing Cloudflare Workers",
            },
          ],
        },
        {
          type: "category",
          label: "Performance Testing",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-playwright-performance-tests-with-artillery-and-tracetest",
              label: "Performance Testing with Distributed Tracing using Artillery, Playwright and Tracetest",
            },
          ],
        },
        {
          type: "category",
          label: "API Gateways",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/testing-distributed-services-with-tyk-opentelemetry-tracetest",
              label: "Testing Distributed Services with Tyk, OpenTelemetry, and Tracetest",
            },
          ],
        },
        {
          type: "category",
          label: "OpenTelemetry",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-without-a-trace-data-store",
              label: "Node.js and OpenTelemetry",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-without-a-trace-data-store-with-manual-instrumentation",
              label: "Node.js and OpenTelemetry Manual Instrumentation",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-python-app-with-opentelemetry-collector-and-tracetest",
              label: "Python and OpenTelemetry Manual Instrumentation",
            },
          ],
        },
        {
          type: "category",
          label: "OpenTelemetry Collector + Tracing Vendors",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-lightstep",
              label: "OpenTelemetry Demo and Lightstep",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-new-relic",
              label: "OpenTelemetry Demo and New Relic",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-datadog",
              label: "OpenTelemetry Demo and Datadog",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-dynatrace",
              label: "OpenTelemetry Demo and Dynatrace",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-honeycomb",
              label: "Node.js and Honeycomb",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-signoz-pokeshop",
              label: "Pokeshop API and SigNoz",
            },
          ],
        },
        {
          type: "category",
          label: "AWS X-Ray",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-aws-x-ray",
              label: "Node.js and AWS X-Ray (Node.js SDK)",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-aws-x-ray-adot",
              label: "Node.js with AWS X-Ray (Node.js SDK) and AWS Distro for OpenTelemetry",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-aws-x-ray-pokeshop",
              label: "Pokeshop API with AWS X-Ray (Node.js SDK) and AWS Distro for OpenTelemetry",
            },

            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-step-functions-terraform",
              label: ".NET Step Functions with AWS X-Ray, AWS Distro for OpenTelemetry, and Terraform",
            },
          ],
        },
        {
          type: "category",
          label: "Azure App Insights",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-azure-app-insights",
              label: "Node.js and Azure Application Insights (Node.js SDK)",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-azure-app-insights-collector",
              label: "Node.js and Azure Application Insights with OpenTelemetry Collector",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-azure-app-insights-pokeshop",
              label: "Pokeshop API and Azure Application Insights with OpenTelemetry Collector",
            },
          ],
        },
        {
          type: "category",
          label: "Elastic",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-elasticapm",
              label: "Node.js and Elastic APM",
            },
          ],
        },
        {
          type: "category",
          label: "Grafana Tempo",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-tempo",
              label: "Node.js and Grafana Tempo",
            },
          ],
        },
        {
          type: "category",
          label: "Jaeger",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-jaeger",
              label: "Node.js and Jaeger",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-aws-terraform",
              label: "Serverless Node.js and Jaeger with Terraform",
            },
          ],
        },
        {
          type: "category",
          label: "OpenSearch",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-opensearch",
              label: "Node.js and OpenSearch",
            },
          ],
        },
        {
          type: "category",
          label: "Sumo Logic",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-sumologic",
              label: "Node.js and Sumo Logic",
            },
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Tools & Integrations",
      link: {
        type: "doc",
        id: "tools-and-integrations/overview",
      },
      items: [
        {
          type: "doc",
          id: "tools-and-integrations/cypress",
          label: "Cypress",
        },
        {
          type: "doc",
          id: "tools-and-integrations/playwright",
          label: "Playwright",
        },
        {
          type: "doc",
          id: "tools-and-integrations/typescript",
          label: "TypeScript",
        },
        {
          type: "doc",
          id: "tools-and-integrations/keptn",
          label: "Keptn",
        },
        {
          type: "doc",
          id: "tools-and-integrations/artillery-plugin",
          label: "Artillery Plugin",
        },
        {
          type: "doc",
          id: "tools-and-integrations/artillery-engine",
          label: "Artillery Engine",
        },
        {
          type: "doc",
          id: "tools-and-integrations/k6",
          label: "k6 & Tracetest",
        },
        {
          type: "doc",
          id: "tools-and-integrations/k6-core",
          label: "k6 & Tracetest Core",
        },
        {
          type: "doc",
          id: "tools-and-integrations/testkube",
          label: "Testkube",
        },
      ],
    },
    {
      type: "category",
      label: "CI/CD Automation",
      link: {
        type: "doc",
        id: "ci-cd-automation/overview",
      },
      items: [
        {
          type: "doc",
          id: "ci-cd-automation/github-actions-pipeline",
          label: "GitHub Actions Pipeline",
        },
        {
          type: "doc",
          id: "ci-cd-automation/testkube-pipeline",
          label: "Testkube Kubernetes-native Test Runner Pipeline",
        },
        {
          type: "doc",
          id: "ci-cd-automation/tekton-pipeline",
          label: "Tekton Cloud-native Pipeline",
        },
      ],
    },
    {
      type: "doc",
      id: "examples-tutorials/videos",
      label: "Videos & Recordings",
    },
    {
      type: "doc",
      id: "examples-tutorials/webinars",
      label: "Webinars",
    },
    {
      type: "doc",
      id: "examples-tutorials/tutorials",
      label: "Blog Post Tutorials",
    },
  ],

  coreSidebar: [
    {
      type: "link",
      label: "Welcome! 👋",
      href: "/",
    },
    {
      type: "category",
      label: "Getting Started",
      link: {
        type: "doc",
        id: "core/getting-started/overview",
      },
      items: [
        {
          type: "doc",
          id: "core/getting-started/installation",
          label: "Installing Tracetest Core",
        },
        {
          type: "doc",
          id: "core/getting-started/open",
          label: "Opening Tracetest Core",
        },
      ],
    },
    {
      type: "category",
      label: "Configuration",
      link: {
        type: "doc",
        id: "core/configuration/overview",
      },
      items: [
        {
          type: "doc",
          id: "core/configuration/server",
          label: "Tracetest Server Configuration",
        },
        {
          type: "doc",
          id: "core/configuration/provisioning",
          label: "Tracetest Server Provisioning",
        },
        {
          type: "doc",
          id: "core/configuration/analytics",
          label: "Tracetest Server Analytics",
        },
        {
          type: "doc",
          id: "core/configuration/telemetry",
          label: "Tracetest Server Telemetry",
        },
        {
          type: "doc",
          id: "core/configuration/upgrade",
          label: "Upgrade Tracetest Version",
        },
      ],
    },
    {
      type: "category",
      label: "Deployment",
      link: {
        type: "doc",
        id: "core/deployment/overview",
      },
      items: [
        {
          type: "doc",
          id: "core/deployment/docker",
          label: "Docker",
        },
        {
          type: "doc",
          id: "core/deployment/kubernetes",
          label: "Kubernetes",
        },
      ],
    },
  ],

  tutorialSidebar: [
    {
      type: "doc",
      id: "index",
      label: "Welcome! 👋",
    },
    {
      type: "category",
      label: "Getting Started",
      link: {
        type: "doc",
        id: "getting-started/overview",
      },
      items: [
        {
          type: "doc",
          id: "getting-started/installation",
          label: "Installing Tracetest",
        },
        {
          type: "doc",
          id: "getting-started/open",
          label: "Opening Tracetest",
        },
        {
          type: "doc",
          id: "getting-started/no-otel",
          label: "What if I don't have OpenTelemetry installed?",
        },
        {
          type: "doc",
          id: "getting-started/try-live-demos",
          label: "Try live demos!",
        },
      ],
    },
    {
      type: "category",
      label: "Configuration",
      link: {
        type: "doc",
        id: "configuration/overview",
      },
      items: [
        {
          type: "doc",
          id: "configuration/agent",
          label: "Tracetest Agent (Private Access)",
        },
        {
          type: "doc",
          id: "configuration/cloud-agent",
          label: "Tracetest Cloud Agent (Public Access)",
        },
        {
          type: "category",
          label: "Connecting to Tracing Backends",
          link: {
            type: "doc",
            id: "configuration/connecting-to-data-stores/overview",
          },
          items: [
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/awsxray",
              label: "AWS X-Ray",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/azure-app-insights",
              label: "Azure App Insights",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/datadog",
              label: "Datadog",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/dynatrace",
              label: "Dynatrace",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/elasticapm",
              label: "Elastic APM",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/tempo",
              label: "Grafana Tempo",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/honeycomb",
              label: "Honeycomb",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/jaeger",
              label: "Jaeger",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/lightstep",
              label: "Lightstep",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/new-relic",
              label: "New Relic",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/opensearch",
              label: "OpenSearch",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/opentelemetry-collector",
              label: "OpenTelemetry Collector",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/signalfx",
              label: "SignalFX",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/signoz",
              label: "SigNoz",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/sumologic",
              label: "Sumo Logic",
            },
            {
              type: "doc",
              id: "configuration/opentelemetry-collector-configuration-file-reference",
              label: "OpenTelemetry Collector Configuration File Reference",
            },
          ],
        },
        {
          type: "doc",
          id: "configuration/tracetest-analyzer",
          label: "Tracetest Analyzer",
        },
        {
          type: "doc",
          id: "configuration/trace-polling",
          label: "Trace Polling",
        },
        {
          type: "doc",
          id: "configuration/test-runner",
          label: "Test Runner",
        },
        {
          type: "doc",
          id: "configuration/demo",
          label: "Demo Applications",
        },
        {
          type: "doc",
          id: "configuration/sampling-tracetest-spans",
          label: "Sampling Tracetest Spans",
        },
      ],
    },
    {
      type: "category",
      label: "Concepts",
      items: [
        {
          type: "doc",
          id: "concepts/what-is-trace-based-testing",
          label: "What is Trace-based Testing",
        },
        // {
        //   type: "doc",
        //   id: "concepts/what-is-tracing",
        //   label: "What is tracing",
        // },
        {
          type: "doc",
          id: "concepts/agent",
          label: "Tracetest Agent (Private Access)",
        },
        {
          type: "doc",
          id: "concepts/cloud-agent",
          label: "Tracetest Cloud Agent (Public Access)",
        },
        {
          type: "doc",
          id: "concepts/polling-profiles",
          label: "Polling Profiles",
        },
        {
          type: "doc",
          id: "concepts/architecture",
          label: "Architecture",
        },
        {
          type: "doc",
          id: "concepts/organizations",
          label: "Organizations",
        },
        {
          type: "doc",
          id: "concepts/environments",
          label: "Environments",
        },
        {
          type: "doc",
          id: "concepts/roles-and-permissions",
          label: "Roles and Permissions",
        },
        {
          type: "doc",
          id: "concepts/environment-tokens",
          label: "Environment Tokens",
        },
        {
          type: "doc",
          id: "concepts/assertions",
          label: "Assertions",
        },
        // {
        //   type: "doc",
        //   id: "concepts/data-stores",
        //   label: "Data Stores",
        // },
        {
          type: "doc",
          id: "concepts/variable-sets",
          label: "Variable Sets",
        },
        {
          type: "doc",
          id: "concepts/selectors",
          label: "Selectors",
        },
        // {
        //   type: "doc",
        //   id: "concepts/tests",
        //   label: "Tests",
        // },
        {
          type: "doc",
          id: "concepts/expressions",
          label: "Expressions",
        },
        {
          type: "doc",
          id: "concepts/test-suites",
          label: "Test Suites",
        },
        {
          type: "doc",
          id: "concepts/ad-hoc-testing",
          label: "Ad-hoc Testing",
        },
        {
          type: "doc",
          id: "concepts/versioning",
          label: "Versioning",
        },
      ],
    },
    {
      type: "category",
      label: "Trace Analyzer",
      items: [
        {
          type: "doc",
          id: "analyzer/concepts",
          label: "Concepts",
        },
        {
          type: "category",
          label: "Plugins",
          items: [
            {
              type: "category",
              label: "OTel Semantic Conventions",
              link: {
                type: "doc",
                id: "analyzer/plugins/otel-semantic-conventions",
              },
              items: [
                {
                  type: "doc",
                  id: "analyzer/rules/span-naming",
                  label: "span-naming",
                },
                {
                  type: "doc",
                  id: "analyzer/rules/attribute-naming",
                  label: "attribute-naming",
                },
                {
                  type: "doc",
                  id: "analyzer/rules/required-attributes",
                  label: "required-attributes",
                },
                {
                  type: "doc",
                  id: "analyzer/rules/no-empty-attributes",
                  label: "no-empty-attributes",
                },
              ],
            },
            {
              type: "category",
              label: "Security",
              link: {
                type: "doc",
                id: "analyzer/plugins/security",
              },
              items: [
                {
                  type: "doc",
                  id: "analyzer/rules/secure-https-protocol",
                  label: "secure-https-protocol",
                },
                {
                  type: "doc",
                  id: "analyzer/rules/no-api-key-leak",
                  label: "no-api-key-leak",
                },
              ],
            },
            {
              type: "category",
              label: "Common Problems",
              link: {
                type: "doc",
                id: "analyzer/plugins/common-problems",
              },
              items: [
                {
                  type: "doc",
                  id: "analyzer/rules/prefer-dns",
                  label: "prefer-dns",
                },
              ],
            },
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Web UI",
      items: [
        {
          type: "doc",
          id: "web-ui/creating-data-stores",
          label: "Configuring Data Stores",
        },
        {
          type: "category",
          label: "Configuring Tests",
          items: [
            {
              type: "category",
              label: "Creating Tests",
              link: {
                type: "doc",
                id: "web-ui/creating-tests",
              },
              items: [
                {
                  type: "doc",
                  id: "web-ui/creating-tests-http",
                  label: "Creating Tests - HTTP",
                },
                {
                  type: "doc",
                  id: "web-ui/creating-tests-gRPC",
                  label: "Creating Tests - gRPC",
                },
                {
                  type: "doc",
                  id: "web-ui/creating-tests-kafka",
                  label: "Creating Tests - Kafka",
                },
                {
                  type: "doc",
                  id: "web-ui/creating-tests-traceid",
                  label: "Creating Tests - TraceID",
                },
                {
                  type: "doc",
                  id: "web-ui/creating-tests-cypress",
                  label: "Creating Tests - Cypress",
                },
                {
                  type: "doc",
                  id: "web-ui/creating-tests-playwright",
                  label: "Creating Tests - Playwright",
                },
                {
                  type: "doc",
                  id: "web-ui/creating-tests-k6",
                  label: "Creating Tests - k6",
                },
                {
                  type: "doc",
                  id: "web-ui/creating-tests-artillery",
                  label: "Creating Tests - Artillery",
                },
              ],
            },
            {
              type: "category",
              label: "Importing Tests",
              items: [
                {
                  type: "doc",
                  id: "web-ui/importing-tests-definition",
                  label: "Importing Tests from Tracetest Definition",
                },
                {
                  type: "doc",
                  id: "web-ui/importing-tests-curl",
                  label: "Importing Tests from CURL",
                },
                {
                  type: "doc",
                  id: "web-ui/importing-tests-postman",
                  label: "Importing Tests from Postman",
                },
              ],
            },
            {
              type: "doc",
              id: "web-ui/creating-test-specifications",
              label: "Creating Test Specifications",
            },
            {
              type: "doc",
              id: "web-ui/creating-test-outputs",
              label: "Creating Test Outputs",
            },
            {
              type: "doc",
              id: "web-ui/test-results",
              label: "Test Results",
            },
            {
              type: "doc",
              id: "web-ui/exporting-tests",
              label: "Exporting Tests",
            },
            {
              type: "doc",
              id: "web-ui/undefined-variables",
              label: "Using Undefined Variables",
            },
          ],
        },
        {
          type: "doc",
          id: "web-ui/creating-test-suites",
          label: "Creating Test Suites",
        },
        {
          type: "doc",
          id: "web-ui/black-box-testing",
          label: "Enabling Black-Box Testing",
        },
        {
          type: "doc",
          id: "web-ui/creating-variable-sets",
          label: "Creating Variable Sets",
        },
      ],
    },
    {
      type: "category",
      label: "CLI",
      items: [
        {
          type: "doc",
          id: "cli/cli-installation-reference",
          label: "CLI Installation Reference",
        },
        {
          type: "doc",
          id: "cli/configuring-your-cli",
          label: "Configuring your CLI",
        },
        {
          type: "doc",
          id: "cli/creating-data-stores",
          label: "Configuring Data Stores",
        },
        {
          type: "category",
          label: "Configuring Tests",
          items: [
            {
              type: "category",
              link: {
                type: "doc",
                id: "cli/creating-tests",
              },
              label: "Creating Tests",
              items: [
                {
                  type: "doc",
                  id: "cli/creating-tests-http",
                  label: "Creating Tests - HTTP",
                },
                {
                  type: "doc",
                  id: "cli/creating-tests-grpc",
                  label: "Creating Tests - gRPC",
                },
                {
                  type: "doc",
                  id: "cli/creating-tests-kafka",
                  label: "Creating Tests - Kafka",
                },
                {
                  type: "doc",
                  id: "cli/creating-tests-traceid",
                  label: "Creating Tests - TraceID",
                },
              ],
            },
            {
              type: "doc",
              id: "cli/creating-test-specifications",
              label: "Creating Test Specifications",
            },
            {
              type: "doc",
              id: "cli/creating-test-outputs",
              label: "Creating Test Outputs",
            },
            {
              type: "doc",
              id: "cli/running-tests",
              label: "Running Tests",
            },
            {
              type: "doc",
              id: "cli/undefined-variables",
              label: "Using Undefined Variables",
            },
            // {
            //   type: "doc",
            //   id: "cli/exporting-tests",
            //   label: "Exporting tests",
            // },
          ],
        },
        {
          type: "category",
          label: "Configuring Test Suites",
          items: [
            {
              type: "doc",
              id: "cli/creating-test-suites",
              label: "Creating Test Suites",
            },
            {
              type: "doc",
              id: "cli/running-test-suites",
              label: "Running Test Suites",
            },
          ],
        },
        {
          type: "doc",
          id: "cli/creating-variable-sets",
          label: "Creating Variable Sets",
        },
        {
          type: "category",
          label: "Reference",
          items: cliGeneratedPages,
        },
      ],
    },
    {
      type: "link",
      label: "CI/CD Automation",
      href: "/ci-cd-automation/overview",
    },
    {
      type: "link",
      label: "Tools & Integrations",
      href: "/tools-and-integrations/overview",
    },
    {
      type: "link",
      label: "Examples & Tutorials",
      href: "/examples-tutorials/overview",
    },
    {
      type: "link",
      label: "Live Examples",
      href: "/live-examples/opentelemetry-store/overview",
    },
    {
      type: "link",
      label: "Tracetest Open API Definition",
      href: "/openapi",
    },
  ],
};

module.exports = sidebars;
