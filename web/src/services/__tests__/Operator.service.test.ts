import {CompareOperator} from '../../constants/Operator.constants';
import OperatorService from '../Operator.service';

describe('OperatorService', () => {
  it('should return the operator name', () => {
    const name = OperatorService.getOperatorName(CompareOperator.EQUALS);

    expect(name).toEqual('equals');
  });

  it('should return the operator symbol', () => {
    const symbol = OperatorService.getOperatorSymbol(CompareOperator.EQUALS);

    expect(symbol).toEqual('==');
  });
});
