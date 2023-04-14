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

export const SupportedDemosFormFieldMap = {
  [SupportedDemosFormField.Pokeshop]: SupportedDemos.Pokeshop,
  [SupportedDemosFormField.OpentelemetryStore]: SupportedDemos.OpentelemetryStore,
};

export enum ResourceType {
  ConfigType = 'Config',
  PollingProfileType = 'PollingProfile',
  DemoType = 'Demo',
}

export enum ResourceTypePlural {
  ConfigType = 'Configs',
  PollingProfileType = 'PollingProfiles',
  DemoType = 'Demos',
}

export type TDraftDemo = Record<Required<Demo['type']>, Partial<Demo>>;
export type TDraftPollingProfiles = Partial<Polling>;
export type TDraftConfig = Partial<Config>;
export type TDraftSpec = TDraftConfig | TDraftPollingProfiles | Partial<Demo>;

export type TDraftResource = {
  type: ResourceType;
  typePlural: ResourceTypePlural;
  spec: TDraftSpec;
};
