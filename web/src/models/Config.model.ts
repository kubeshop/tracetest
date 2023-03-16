import {Model, TConfigSchemas} from 'types/Common.types';

export type TRawConfig = TConfigSchemas['ConfigurationResource'];

type Config = Model<Model<TRawConfig, {}>['spec'], {}>;

function Config({
  spec: {analyticsEnabled = false, id = 'current', name = 'Config'} = {
    analyticsEnabled: false,
  },
}: TRawConfig = {}): Config {
  return {
    analyticsEnabled,
    id,
    name,
  };
}

export type TRawLiveConfig = {
  AnalyticsEnabled?: boolean;
  ID: 'current';
  Name: 'Config';
};

Config.FromLiveUpdate = ({AnalyticsEnabled: analyticsEnabled = false, Name: name, ID: id}: TRawLiveConfig): Config => {
  return {
    id,
    name,
    analyticsEnabled,
  };
};

export default Config;
