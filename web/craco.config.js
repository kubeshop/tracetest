const CracoLessPlugin = require('craco-less');

module.exports = {
  plugins: [
    {
      plugin: CracoLessPlugin,
      options: {
        lessLoaderOptions: {
          lessOptions: {
            modifyVars: {
              /** General Color */
              '@primary-color': '#61175E',
              '@processing-color': '#61175E',
              '@heading-color': '#031849',
              '@text-color': '#031849',
              '@text-color-secondary': '#687492',
              /** Tooltip */
              '@tooltip-color': '#031849',
              '@tooltip-bg': '#FBFBFF',
              /** Select */
              '@select-background': '#FAFAFA',
              '@select-item-selected-bg': '#FAFAFA',
              /** Heading */
              '@heading-1-size': '18px',
              '@heading-2-size': '16px',
              '@heading-3-size': '14px',
              '@heading-4-size': '12px',
              '@heading-5-size': '10px',
            },
            javascriptEnabled: true,
          },
        },
      },
    },
  ],
};
