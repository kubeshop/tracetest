import {PseudoSelector} from 'constants/Operator.constants';
import AssertionSpanResult from '../models/AssertionSpanResult.model';
import {Model} from './Common.types';
import {TCompareOperatorSymbol} from './Operator.types';
import {TSpanFlatAttribute} from './Span.types';

export type TSpanSelector = Model<
  TSpanFlatAttribute,
  {
    operator: TCompareOperatorSymbol;
  }
>;

export type TPseudoSelector = {
  selector: PseudoSelector;
  number?: number;
};

export type TStructuredAssertion = {
  left: string;
  comparator: TCompareOperatorSymbol;
  right: string;
};

export interface ICheckResult {
  result: AssertionSpanResult;
  assertion: string;
}
