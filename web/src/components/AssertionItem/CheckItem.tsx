import AttributeValue from 'components/AttributeValue';
import OperatorService from 'services/Operator.service';
import {ICheckResult} from 'types/Assertion.types';
import {TCompareOperatorSymbol} from 'types/Operator.types';
import * as S from './AssertionItem.styled';

interface IProps {
  check: ICheckResult;
}

const CheckItem = ({check}: IProps) => (
  <S.GridContainer>
    <S.Row>
      {check.result.passed ? <S.IconSuccess /> : <S.IconError />}
      <S.CheckContainer>
        {check.assertion.attribute}
        {` `}
        <S.SecondaryText>
          {OperatorService.getNameFromSymbol(check.assertion.comparator as TCompareOperatorSymbol)}
        </S.SecondaryText>
        {` `}
        <AttributeValue value={check.assertion.expected} />
      </S.CheckContainer>
    </S.Row>

    <S.Row $align="end">
      <AttributeValue
        strong
        type={check.result.passed ? 'success' : 'danger'}
        value={check.result.observedValue || '<Empty Value>'}
      />
    </S.Row>
  </S.GridContainer>
);

export default CheckItem;
