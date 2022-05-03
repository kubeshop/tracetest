import {CompareOperatorNameMap, CompareOperatorSymbolMap} from '../constants/Operator.constants';
import {TCompareOperator, TCompareOperatorName, TCompareOperatorSymbol} from '../types/Operator.types';

const OperatorService = () => ({
  getOperatorName(op: TCompareOperator): TCompareOperatorName {
    return CompareOperatorNameMap[op];
  },
  getOperatorSymbol(op: TCompareOperator): TCompareOperatorSymbol {
    return CompareOperatorSymbolMap[op];
  },
});

export default OperatorService();
