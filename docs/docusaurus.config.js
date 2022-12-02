// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Tracetest Docs',
  tagline: 'Trace-based testing',
  url: 'https://tracetest.io',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/logo.svg',

  // GitHub pages deployment config.
  // If you aren't using GitHub pages, you don't need these.
  organizationName: 'kubeshop', // Usually your GitHub org/user name.
  projectName: 'tracetest', // Usually your repo name.

  // Even if you don't use internalization, you can use this field to set useful
  // metadata like html lang. For example, if your site is Chinese, you may want
  // to replace "en" with "zh-Hans".
  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  markdown: {
    mermaid: true,
  },
  themes: ['@docusaurus/theme-mermaid'],

  presets: [
    [
      "classic",
      /** @type {import('@docusaurus/preset-classic').Options} */
      ({
        docs: {
          sidebarPath: require.resolve("./sidebars.js"),
          // Please change this to your repo.
          // Remove this to remove the "edit this page" links.
          editUrl: "https://github.com/kubeshop/tracetest/blob/main/docs/",
          routeBasePath: "/",
        },
        theme: {
          customCss: require.resolve("./src/css/custom.css"),
        },
        sitemap: {
          changefreq: 'always',
          priority: 0.5,
          ignorePatterns: ['/tags/**'],
          filename: 'sitemap.xml',
        },
      }),
    ],
    [
      "redocusaurus",
      {
        // Plugin Options for loading OpenAPI files
        specs: [
          {
            spec: "https://raw.githubusercontent.com/kubeshop/tracetest/main/api/openapi.yaml",
            route: "/openapi/",
          },
        ],
        // Theme Options for modifying how redoc renders them
        theme: {
          // Change with your site colors
          primaryColor: "#1890ff",
        },
      },
    ],
  ],
  plugins: [
    [
      require.resolve('docusaurus-gtm-plugin'),
      {
        id: 'GTM-MZ7RNS7', // GTM Container ID
      },
    ]
  ],
  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      colorMode: {
        defaultMode: 'dark',
        disableSwitch: false,
        respectPrefersColorScheme: false,
      },  
      navbar: {
        title: 'Tracetest',
        logo: {
          alt: 'Tracetest Logo',
          src: 'img/logo.svg',
        },
        items: [
          // {
          //   to: "/quick-start",
          //   label: "Quick Start",
          //   position: "left",
          // },
          {
            href: "https://discord.gg/6zupCZFQbe",
            label: "Discord",
            position: "left",
          },
          {
            type: "html",
            position: "right",
            value: `<iframe src="https://ghbtns.com/github-btn.html?user=kubeshop&repo=tracetest&type=star&count=true&size=large" style='margin-top: 6px' frameborder="0" scrolling="0" width="170" height="30" title="GitHub"></iframe>`,
          },
          {
            type: "search",
            position: "left",
          },
        ],
      },
      footer: {
        style: "dark",
        links: [
          {
            title: "Developers",
            items: [
              {
                label: "Docs",
                to: "/",
              },
            ],
          },
          {
            title: "Community",
            items: [
              {
                label: "Discord",
                href: "https://discord.gg/6zupCZFQbe",
              },
              {
                label: "Twitter",
                href: "https://twitter.com/tracetest_io",
              },
            ],
          },
          {
            title: "More",
            items: [
              {
                label: "Blog",
                to: "https://tracetest.io/blog",
              },
              {
                label: "GitHub",
                href: "https://github.com/kubeshop/tracetest",
              },
            ],
          },
        ],
        copyright: `Copyright © ${new Date().getFullYear()} Kubeshop, Inc.`,
      },
      prism: {
        theme: lightCodeTheme,
        darkTheme: darkCodeTheme,
      },
      algolia: {
        // The application ID provided by Algolia
        appId: "NA2OEQH0PY",

        // Public API key: it is safe to commit it
        apiKey: "4bbc30616d55bf1aad84ed654121d5c7",

        indexName: "tracetest",
      },
    }),
};

module.exports = config;
