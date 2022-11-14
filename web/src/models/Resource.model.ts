import {ResourceType, TRawResource, TResource} from 'types/Resource.type';
import {TRawTest} from 'types/Test.types';
import {TRawTransaction} from 'types/Transaction.types';
import Test from './Test.model';
import Transaction from './Transaction.model';

function Resource({item, type}: TRawResource): TResource {
  if (type === ResourceType.Test) {
    return {
      type: ResourceType.Test,
      item: Test(item as TRawTest),
    };
  }

  return {
    type: ResourceType.Transaction,
    item: Transaction(item as TRawTransaction),
  };
}

export default Resource;
