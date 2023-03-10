import {TRawConfig} from 'models/Config.model';
import {TRawDemo} from 'models/Demo.model';
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
  Demo = 'Demo',
  Polling = 'PollingProfile',
}

export type TSpec = TRawConfig | TRawPolling | TRawDemo;
