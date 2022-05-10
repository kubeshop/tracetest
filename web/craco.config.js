const CracoLessPlugin = require('craco-less');

module.exports = {
  plugins: [
    {
      plugin: CracoLessPlugin,
      options: {
        lessLoaderOptions: {
          lessOptions: {
            modifyVars: {
              '@primary-color': '#61175E',
              '@border-radius-base': '2px',
            },
            javascriptEnabled: true,
          },
        },
      },
    },
  ],
};
