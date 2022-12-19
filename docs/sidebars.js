/**
 * Creating a sidebar enables you to:
 - create an ordered group of docs
 - render a sidebar for each doc of that group
 - provide next/previous navigation

 The sidebars can be generated from the filesystem, or explicitly defined here.

 Create as many sidebars as you want.
 */

// @ts-check

/** @type {import('@docusaurus/plugin-content-docs').SidebarsConfig} */
const sidebars = {
  // By default, Docusaurus generates a sidebar from the docs folder structure
  // tutorialSidebar: [{type: 'autogenerated', dirName: '.'}],

  // But you can create a sidebar manually

  tutorialSidebar: [
    {
      type: "doc",
      id: "index",
      label: "Introduction",
    },
    {
      type: "category",
      label: "Getting Started",
      items: [
        {
          type: "doc",
          id: "getting-started/installation",
          label: "Quick Start",
        },
        {
          type: "doc",
          id: "getting-started/detailed-installation",
          label: "Detailed Installation",
        },
      ],
    },
    {
      type: "category",
      label: "Configuration",
      items: [
        {
          type: "doc",
          id: "configuration/overview",
          label: "Overview",
        },
        // {
        //   type: "doc",
        //   id: "configuration/config-file-reference",
        //   label: "Config File Reference",
        // },
        {
          type: "category",
          label: "Connecting to Data Stores",
          items: [
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/opentelemetry-collector",
              label: "OpenTelemetry Collector",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/jaeger",
              label: "Jaeger",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/opensearch",
              label: "OpenSearch",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/signalfx",
              label: "SignalFX",
            },
            {
              type: "doc",
              id: "configuration/connecting-to-data-stores/tempo",
              label: "Tempo",
            },
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Deployment",
      items: [
        {
          type: "doc",
          id: "deployment/overview",
          label: "Deployment Overview",
        },
        // {
        //   type: "doc",
        //   id: "deployment/production-checklist",
        //   label: "Production checklist",
        // },
        {
          type: "doc",
          id: "deployment/docker",
          label: "Docker",
        },
        {
          type: "doc",
          id: "deployment/kubernetes",
          label: "Kubernetes",
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
          label: "What is trace-based testing",
        },
        // {
        //   type: "doc",
        //   id: "concepts/what-is-tracing",
        //   label: "What is tracing",
        // },
        {
          type: "doc",
          id: "concepts/architecture",
          label: "Architecture",
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
          id: "concepts/environments",
          label: "Environments",
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
          id: "concepts/transactions",
          label: "Transactions",
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
      label: "Web UI",
      items: [
        // {
        //   type: "doc",
        //   id: "web-ui/creating-environments",
        //   label: "Creating environments",
        // },
        {
          type: "doc",
          id: "web-ui/creating-tests",
          label: "Creating tests",
        },
        {
          type: "doc",
          id: "web-ui/creating-test-specifications",
          label: "Creating test specifications",
        },
        {
          type: "doc",
          id: "web-ui/test-results",
          label: "Test results",
        },
        // {
        //   type: "doc",
        //   id: "web-ui/creating-transactions",
        //   label: "Creating transactions",
        // },
        {
          type: "doc",
          id: "web-ui/exporting-tests",
          label: "Exporting tests",
        },
      ],
    },
    {
      type: "category",
      label: "CLI",
      items: [
        {
          type: "doc",
          id: "cli/configuring-your-cli",
          label: "Configuring your CLI",
        },
        // {
        //   type: "doc",
        //   id: "cli/creating-environments",
        //   label: "Creating environments",
        // },
        {
          type: "doc",
          id: "cli/creating-tests",
          label: "Creating tests",
        },
        // {
        //   type: "doc",
        //   id: "cli/creating-transactions",
        //   label: "Creating transactions",
        // },
        // {
        //   type: "doc",
        //   id: "cli/exporting-tests",
        //   label: "Exporting tests",
        // },
      ],
    },
    {
      type: "category",
      label: "CI/CD automation",
      items: [
        {
          type: "doc",
          id: "ci-cd-automation/github-actions-pipeline",
          label: "GitHub actions pipeline",
        },
      ],
    },
    {
      type: "category",
      label: "Examples & Tutorials",
      items: [
        {
          type: "doc",
          id: "examples-tutorials/overview",
          label: "Overview",
        },
        {
          type: "category",
          label: "Recipes",
          items: [
            {
              type: "doc",
              id: "examples-tutorials/recipes",
              label: "Overview",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-without-a-trace-data-store",
              label: "Running Tracetest Without a Trace Data Store",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-jaeger",
              label: "Running Tracetest With Jaeger",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-opensearch",
              label: "Running Tracetest With OpenSearch",
            },
            {
              type: "doc",
              id: "examples-tutorials/recipes/running-tracetest-with-tempo",
              label: "Running Tracetest With Tempo",
            },
          ],
        },
      ],
    },
    {
      type: "category",
      label: "Live examples",
      items: [
        {
          type: "category",
          label: "Pokemon API Demo",
          items: [
            {
              type: "doc",
              id: "pokeshop/index",
              label: "Overview",
            },
            {
              type: "doc",
              id: "pokeshop/add-pokemon",
              label: "Add Pokemon",
            },
            {
              type: "doc",
              id: "pokeshop/import-pokemon",
              label: "Import Pokemon",
            },
            {
              type: "doc",
              id: "pokeshop/list-pokemon",
              label: "List Pokemon",
            },
          ],
        },
        // {
        //   type: "category",
        //   label: "OpenTelemetry Store Demo",
        //   items: [
        //     {
        //       type: "doc",
        //       id: "opentelemetry-store/overview",
        //       label: "Overview",
        //     },
        //   ],
        // },
      ],
    },
    {
      type: "link",
      label: "Tracetest Open API definition",
      href: "/openapi",
    },
  ],
};

module.exports = sidebars;