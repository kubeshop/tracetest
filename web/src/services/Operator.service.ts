import {CompareOperator, CompareOperatorNameMap, CompareOperatorSymbolMap} from '../constants/Operator.constants';
import {TCompareOperatorName, TCompareOperatorSymbol} from '../types/Operator.types';

const OperatorService = () => ({
  getOperatorName(op: CompareOperator): TCompareOperatorName {
    return CompareOperatorNameMap[op];
  },
  getOperatorSymbol(op: CompareOperator): TCompareOperatorSymbol {
    return CompareOperatorSymbolMap[op];
  },
});

export default OperatorService();
