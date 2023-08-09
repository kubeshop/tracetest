import {ResourceType} from 'types/Resource.type';
import {Model, TResourceSchemas} from '../types/Common.types';

import Test, {TRawTest} from './Test.model';
import TestSuite, {TRawTestSuite} from './TestSuite.model';

export type TRawResource = TResourceSchemas['Resource'];
type Resource = Model<TRawResource, {type: ResourceType; item: Test | TestSuite}>;

function Resource({item, type}: TRawResource): Resource {
  if (type === ResourceType.Test) {
    return {
      type: ResourceType.Test,
      item: Test.FromRawTest(item as TRawTest),
    };
  }

  return {
    type: ResourceType.TestSuite,
    item: TestSuite.FromRawTestSuite(item as TRawTestSuite),
  };
}

export default Resource;
