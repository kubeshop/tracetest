import {ResourceType} from 'types/Resource.type';
import {Model, TResourceSchemas} from '../types/Common.types';

import Test, {TRawTest} from './Test.model';
import Transaction, {TRawTransaction} from './Transaction.model';

export type TRawResource = TResourceSchemas['Resource'];
type Resource = Model<TRawResource, {type: ResourceType; item: Test | Transaction}>;

function Resource({item, type}: TRawResource): Resource {
  if (type === ResourceType.Test) {
    return {
      type: ResourceType.Test,
      item: Test(item as TRawTest),
    };
  }

  return {
    type: ResourceType.Transaction,
    item: Transaction.FromRawTransaction(item as TRawTransaction),
  };
}

export default Resource;
