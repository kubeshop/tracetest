import {IEnv} from 'types/Common.types';

type TEnv = keyof IEnv;

const emptyValues: IEnv = {
  analyticsEnabled: false,
  appVersion: '',
  demoEnabled: [],
  demoEndpoints: {},
  env: '',
  experimentalFeatures: [],
  measurementId: '',
  serverID: '',
  serverPathPrefix: '/',
};

const initialEnv = window.ENV || {};

const Env = {
  get<Key extends TEnv>(key: Key) {
    return initialEnv[key] ?? emptyValues[key];
  },
};

export default Env;
