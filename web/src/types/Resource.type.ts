import {Model, TResourceSchemas} from './Common.types';
import {TTest} from './Test.types';
import {TTransaction} from './Transaction.types';

export enum ResourceType {
  test = 'test',
  transaction = 'transaction',
}

export type TRawResource = TResourceSchemas['Resource'];

export type TResource = Model<TRawResource, {type: ResourceType; item: TTest | TTransaction}>;
