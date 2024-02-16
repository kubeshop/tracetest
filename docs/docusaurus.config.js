// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion

const lightCodeTheme = require('prism-react-renderer').themes.github;
const darkCodeTheme = require('prism-react-renderer').themes.dracula;

/** @type {import('@docusaurus/types').Config} */
const config = {
  title: 'Tracetest Docs',
  tagline: 'Trace-based testing',
  url: 'https://docs.tracetest.io',
  baseUrl: '/',
  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',
  favicon: 'img/logo.svg',
  trailingSlash: false,

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
    // [
    //   '@docusaurus/plugin-google-gtag',
    //   {
    //     trackingID: 'G-999X9XX9XX',
    //     anonymizeIP: true,
    //   },
    // ],
    [
      '@docusaurus/plugin-ideal-image',
      {
        quality: 70,
        max: 1030, // max resized image's size.
        min: 640, // min resized image's size. if original is lower, use that size.
        steps: 2, // the max number of images generated between min and max (inclusive)
        disableInDev: false,
      },
    ],
    [
      require.resolve('docusaurus-gtm-plugin'),
      {
        id: 'GTM-MZ7RNS7', // GTM Container ID
      },
    ],
  ],
  themeConfig:
    /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
    ({
      docs: {
        sidebar: {
          hideable: true,
          autoCollapseCategories: true,
        },
      },  
      colorMode: {
        defaultMode: 'light',
        disableSwitch: false,
        respectPrefersColorScheme: false,
      },
      // Use this to add an announcement for a webinar or event.
      announcementBar: {
        id: 'announcement',
        // content:
        //   '<a target="_blank" rel="noopener noreferrer" href="https://www.youtube.com/live/2MSDy3XHjtE?si=VlK7cxJOsgKi5QTE&t=1132">Tracetest is the official testing harness for the OpenTelemetry Demo! 🚀</a>',
        content:
          '<a target="_blank" rel="noopener noreferrer" href="https://tracetest.io/blog/opentelemetry-is-not-just-for-monitoring-and-troubleshooting-any-longer-announcing-tracetest-open-beta">Tracetest Open Beta is Live. Try it! Give us feedback! 🙌</a>',
        isCloseable: false,
      },  
      navbar: {
        hideOnScroll: true,
        logo: {
          alt: 'Tracetest Logo',
          src: 'img/logo-landscape.svg',
          srcDark: 'img/logo-landscape-dark.svg',
          href: 'https://tracetest.io/',
          target: '_blank',
        },
        items: [
          {
            type: 'docSidebar',
            position: 'left',
            sidebarId: 'tutorialSidebar',
            label: 'Docs'
          },
          {
            type: 'docSidebar',
            position: 'left',
            sidebarId: 'coreSidebar',
            label: 'Core'
          },
          {
            type: 'docSidebar',
            position: 'left',
            sidebarId: 'examplesTutorialsSidebar',
            label: 'Examples & Tutorials'
          },
          {
            type: 'docSidebar',
            position: 'left',
            sidebarId: 'liveExamplesSidebar',
            label: 'Live Examples'
          },
          {
            type: 'dropdown',
            label: 'Support',
            position: 'left',
            items: [
              {
                label: 'Overview',
                href: 'https://tracetest.io/support',
              },
              {
                label: 'Community',
                href: 'https://tracetest.io/community',
              },
              {
                label: 'Pricing',
                href: 'https://tracetest.io/pricing',
              },
              {
                label: 'Talk to us in Slack',
                href: 'https://dub.sh/tracetest-community',
              },
              {
                label: 'Open an issue in GitHub',
                href: 'https://github.com/kubeshop/tracetest/issues/new/choose',
              },
              {
                label: 'Learn',
                href: 'https://tracetest.io/learn',
              },
              {
                label: 'Contact Us',
                href: 'https://tracetest.io/contact',
              },
            ],
          },
          {
            type: "html",
            position: "left",
            value: `<iframe src="https://ghbtns.com/github-btn.html?user=kubeshop&repo=tracetest&type=star&count=true&size=medium" style='margin-top: 6px' frameborder="0" scrolling="0" width="90" height="20" title="GitHub"></iframe>`,
          },
          {
            type: "search",
            position: "right",
          },
          {
            href: "https://app.tracetest.io",
            label: "Sign In",
            position: "right",
            className: "sign-in-button",
          },
        ],
      },
      footer: {
        style: "light",
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
                label: "Slack",
                href: "https://dub.sh/tracetest-community",
              },
              {
                label: "Twitter",
                href: "https://twitter.com/tracetest_io",
              },
              {
                label: "LinkedIn",
                href: "https://www.linkedin.com/company/87135575",
              },
            ],
          },
          {
            title: "More",
            items: [
              {
                label: "Home",
                to: "https://tracetest.io",
              },
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
        appId: "L2ILN3GLIL",

        // Public API key: it is safe to commit it
        apiKey: "663c91299e298ff34c5a7a18f4451d1a",

        indexName: "tracetest",

        contextualSearch: false,

        searchPagePath: false,
      },
    }),
};

module.exports = config;
