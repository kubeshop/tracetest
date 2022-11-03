import {ResourceType, TRawResource, TResource} from 'types/Resource.type';
import {TRawTest} from 'types/Test.types';
import {TRawTransaction} from 'types/Transaction.types';
import Test from './Test.model';
import Transaction from './Transaction.model';

function Resource({item, type}: TRawResource): TResource {
  return {
    type: type === ResourceType.test ? ResourceType.test : ResourceType.transaction,
    item: type === ResourceType.test ? Test(item as TRawTest) : Transaction(item as TRawTransaction),
  };
}

export default Resource;
