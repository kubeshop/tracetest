import Config from 'models/Config.model';
import Polling from 'models/Polling.model';
import Demo from 'models/Demo.model';

export type TListResponse<T> = {
  count: number;
  items: T[];
};

export enum SupportedDemos {
  Pokeshop = 'pokeshop',
  OpentelemetryStore = 'otelstore',
}

export enum SupportedDemosFormField {
  Pokeshop = 'pokeshop',
  OpentelemetryStore = 'opentelemetryStore',
}

export enum ResourceType {
  ConfigType = 'Config',
  PollingProfileType = 'PollingProfile',
  DemoType = 'Demo',
}

export type TDraftDemo = Record<Required<Demo['type']>, Partial<Demo>>;
export type TDraftPollingProfiles = Partial<Polling>;
export type TDraftConfig = Partial<Config>;
export type TDraftSpec = TDraftConfig | TDraftPollingProfiles | Partial<Demo>;

export type TDraftResource = {
  type: ResourceType;
  spec: TDraftSpec;
};
