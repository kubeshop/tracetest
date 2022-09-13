const CracoLessPlugin = require('craco-less');
const webpack = require('webpack');

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
              '@border-color-base': '#CDD1DB',
              /** Tooltip */
              '@tooltip-color': '#031849',
              '@tooltip-bg': '#FBFBFF',
              /** Select */
              '@select-background': '#FFFFFF',
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
  webpack: {
    configure: webpackConfig => {
      return {
        ...webpackConfig,
        module: {
          ...webpackConfig.module,
          rules: [
            ...webpackConfig.module.rules,
            {
              test: /\.m?js/,
              resolve: {
                fullySpecified: false,
              },
            },
          ],
        },
        plugins: [
          ...webpackConfig.plugins,
          new webpack.ProvidePlugin({Buffer: ['buffer', 'Buffer']}),
          new webpack.ProvidePlugin({process: 'process/browser'}),
        ],
        resolve: {
          ...webpackConfig.resolve,
          fallback: {
            stream: require.resolve('stream-browserify'),
            buffer: require.resolve('buffer'),
            path: require.resolve('path-browserify'),
          },
        },
      };
    },
  },
};
