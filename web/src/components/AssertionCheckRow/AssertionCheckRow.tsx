import {useMemo} from 'react';
import {difference} from 'lodash';
import OperatorService from '../../services/Operator.service';
import {ISpanAssertionResult} from '../../types/Assertion.types';
import {ISpan} from '../../types/Span.types';
import * as S from './AssertionCheckRow.styled';
import AttributeValue from '../AttributeValue';

interface IAssertionCheckRowProps {
  result: ISpanAssertionResult;
  span: ISpan;
  assertionSelectorList: string[];
  getIsSelectedSpan(spanId: string): boolean;
  onSelectSpan(spanId: string): void;
}

const AssertionCheckRow: React.FC<IAssertionCheckRowProps> = ({
  result: {propertyName, comparisonValue, operator, actualValue, hasPassed, spanId},
  span: {signature},
  assertionSelectorList,
  getIsSelectedSpan,
  onSelectSpan,
}) => {
  const signatureSelectorList = signature.map(({value}) => value).concat([`#${spanId.slice(-4)}`]) || [];
  const spanLabelList = difference(signatureSelectorList, assertionSelectorList);
  const badgeList = useMemo(() => {
    const isSelected = getIsSelectedSpan(spanId);

    return (isSelected ? [<S.SelectedLabelBadge count="selected" key="selected" />] : []).concat(
      spanLabelList
        // eslint-disable-next-line react/no-array-index-key
        .map((label, index) => <S.LabelBadge count={label} key={`${label}-${index}`} />)
    );
  }, [getIsSelectedSpan, spanId, spanLabelList]);

  return (
    <S.AssertionCheckRow onClick={() => onSelectSpan(spanId)}>
      <S.Entry>
        <S.Label>Span Labels</S.Label>
        <S.Value>{badgeList}</S.Value>
      </S.Entry>
      <S.Entry>
        <S.Label>Attribute</S.Label>
        <S.Value>{propertyName}</S.Value>
      </S.Entry>
      <S.Entry>
        <S.Label>Assertion Type</S.Label>
        <S.Value>{OperatorService.getOperatorName(operator)}</S.Value>
      </S.Entry>
      <S.Entry>
        <S.Label>Expected Value</S.Label>
        <AttributeValue value={comparisonValue} />
      </S.Entry>
      <S.Entry>
        <S.Label>Actual Value</S.Label>
        <AttributeValue strong type={hasPassed ? 'success' : 'danger'} value={actualValue || '<Empty Value>'} />
      </S.Entry>
    </S.AssertionCheckRow>
  );
};

export default AssertionCheckRow;
