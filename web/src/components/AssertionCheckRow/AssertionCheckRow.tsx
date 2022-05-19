import {difference} from 'lodash';
import {useMemo} from 'react';
import OperatorService from '../../services/Operator.service';
import {ISpanAssertionResult} from '../../types/Assertion.types';
import {ISpan} from '../../types/Span.types';
import AttributeValue from '../AttributeValue';
import * as S from './AssertionCheckRow.styled';

interface IAssertionCheckRowProps {
  result: ISpanAssertionResult;
  span: ISpan;
  assertionSelectorList: string[];

  getIsSelectedSpan(spanId: string): boolean;

  onSelectSpan(spanId: string): void;
}

const AssertionCheckRow: React.FC<IAssertionCheckRowProps> = ({
  result: {propertyName, comparisonValue, operator, actualValue, hasPassed, spanId},
  span: {signature, type},
  assertionSelectorList,
  getIsSelectedSpan,
  onSelectSpan,
}) => {
  const signatureSelectorList = signature.map(({value}) => value).concat([`#${spanId.slice(-4)}`]) || [];
  const spanLabelList = difference(signatureSelectorList, assertionSelectorList);
  const badgeList = useMemo(() => {
    const isSelected = getIsSelectedSpan(spanId);

    return (isSelected ? [<S.SelectedLabelBadge spanType={type} count="selected" key="selected" />] : []).concat(
      spanLabelList.map((label, index) => (
        // eslint-disable-next-line react/no-array-index-key
        <S.LabelBadge spanType={index === 0 ? type : undefined} count={label} key={`${label}-${index}`} />
      ))
    );
  }, [getIsSelectedSpan, spanId, spanLabelList, type]);

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
