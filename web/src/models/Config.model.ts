import {EResourceType, TResource} from 'types/Settings.types';

export type TRawConfig = {
  analyticsEnabled: boolean;
  id: string;
  name: string;
};

type Config = {
  analyticsEnabled: boolean;
  id: string;
  name: string;
};

function Config(rawConfig?: TResource<TRawConfig>): Config {
  return {
    analyticsEnabled: rawConfig?.spec?.analyticsEnabled ?? false,
    id: rawConfig?.spec?.id ?? 'current',
    name: rawConfig?.spec?.name ?? '',
  };
}

export function rawToResource(spec: TRawConfig): TResource<TRawConfig> {
  return {spec, type: EResourceType.Config};
}

export default Config;
