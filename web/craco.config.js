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
        plugins: [...webpackConfig.plugins, new webpack.ProvidePlugin({Buffer: ['buffer', 'Buffer']})],
        resolve: {
          ...webpackConfig.resolve,
          fallback: {
            path: require.resolve('path-browserify'),
          },
        },
      };
    },
  },
};
