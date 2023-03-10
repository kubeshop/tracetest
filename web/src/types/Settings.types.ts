import {TRawConfig} from 'models/Config.model';
import {TRawPolling} from 'models/Polling.model';

export type TListResponse<T> = {
  count: number;
  items: TResource<T>[];
};

export type TResource<T> = {
  spec: T;
  type: EResourceType;
};

export enum EResourceType {
  Config = 'Config',
  Polling = 'Polling',
}

export type TSpec = TRawConfig | TRawPolling;

export interface IDraftSettings {
  analytics: boolean;
  maxWaitTimeForTrace: string;
  retryDelay: string;
  demo: IDemoSettings;
}

interface IDemoSettings {
  pokeshopEnabled: boolean;
  pokeshopHttp: string;
  pokeshopGrpc: string;
  otelEnabled: boolean;
  otelFrontend: string;
  otelProductCatalog: string;
  otelCart: string;
  otelCheckout: string;
}
