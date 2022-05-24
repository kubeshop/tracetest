import {
  CompareOperator,
  CompareOperatorNameMap,
  CompareOperatorSymbolMap,
  CompareOperatorSymbolNameMap,
} from '../constants/Operator.constants';
import {TCompareOperatorName, TCompareOperatorSymbol} from '../types/Operator.types';

const OperatorService = () => ({
  getOperatorName(op: CompareOperator): TCompareOperatorName {
    return CompareOperatorNameMap[op];
  },
  getOperatorSymbol(op: CompareOperator): TCompareOperatorSymbol {
    return CompareOperatorSymbolMap[op];
  },
  getNameFromSymbol(symbol: TCompareOperatorSymbol): TCompareOperatorName {
    return CompareOperatorSymbolNameMap[symbol];
  },
});

export default OperatorService();
