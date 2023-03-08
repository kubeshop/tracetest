import {IEnv} from 'types/Common.types';

type TEnv = keyof IEnv;

const emptyValues: IEnv = {
  analyticsEnabled: false,
  appVersion: '',
  env: '',
  experimentalFeatures: [],
  measurementId: '',
  serverID: '',
  serverPathPrefix: '/',
};

const Env = {
  get<Key extends TEnv>(key: Key) {
    const initialEnv = window.ENV || {};
    return initialEnv[key] ?? emptyValues[key];
  },
  set<Key extends TEnv>(key: Key, value: any) {
    window.ENV = {...window.ENV, [key]: value};
  },
};

export default Env;
