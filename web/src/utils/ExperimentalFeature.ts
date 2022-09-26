const {experimentalFeatures = '[]'} = window.ENV || {};
const parsedExperimentalFeatures = JSON.parse(experimentalFeatures);

const ExperimentalFeature = {
  isEnabled(feature: string) {
    return parsedExperimentalFeatures?.includes(feature);
  },
};

export default ExperimentalFeature;
