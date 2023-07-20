const webpack = require('webpack');

module.exports = {
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
            url: require.resolve('url/'),
            fs: false,
          },
        },
      };
    },
  },
};
