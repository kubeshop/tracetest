import {TAssertion} from './Assertion.types';
import {Modify, Schemas} from './Common.types';

export type THTTPRequest = Schemas['HTTPRequest'];
export type TTest = Modify<
  Schemas['Test'],
  {
    assertions: TAssertion[];
  }
>;
