import Env from './Env';

const experimentalFeatures = Env.get('experimentalFeatures');

const ExperimentalFeature = {
  isEnabled(feature: string) {
    return experimentalFeatures.includes(feature);
  },
};

export default ExperimentalFeature;
