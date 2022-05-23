import {useMemo} from 'react';
import {difference} from 'lodash';
import {TAssertion, TAssertionSpanResult} from '../../types/Assertion.types';
import * as S from './AssertionCheckRow.styled';
import AttributeValue from '../AttributeValue';
import {useTestRun} from '../../providers/TestRun/TestRun.provider';

interface TAssertionCheckRowProps {
  result: TAssertionSpanResult;
  assertion: TAssertion;
  assertionSelectorList: string[];

  getIsSelectedSpan(spanId: string): boolean;

  onSelectSpan(spanId: string): void;
}

const AssertionCheckRow: React.FC<TAssertionCheckRowProps> = ({
  result: {observedValue, passed, spanId},
  assertion: {attribute, comparator, expected},
  assertionSelectorList,
  getIsSelectedSpan,
  onSelectSpan,
}) => {
  const {
    run: {trace},
  } = useTestRun();
  const span = useMemo(() => trace?.spans.find(({id}) => id === spanId), [spanId, trace?.spans]);

  const signatureSelectorList = span?.signature.map(({value}) => value).concat([`#${spanId.slice(-4)}`]) || [];
  const spanLabelList = difference(signatureSelectorList, assertionSelectorList);
  const badgeList = useMemo(() => {
    const isSelected = getIsSelectedSpan(spanId);

    return (isSelected ? [<S.SelectedLabelBadge count="selected" key="selected" />] : []).concat(
      spanLabelList.map((label, index) => (
        // eslint-disable-next-line react/no-array-index-key
        <S.LabelBadge count={label} key={`${label}-${index}`} />
      ))
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
        <S.Value>{attribute}</S.Value>
      </S.Entry>
      <S.Entry>
        <S.Label>Assertion Type</S.Label>
        <S.Value>{comparator}</S.Value>
      </S.Entry>
      <S.Entry>
        <S.Label>Expected Value</S.Label>
        <AttributeValue value={expected} />
      </S.Entry>
      <S.Entry>
        <S.Label>Actual Value</S.Label>
        <AttributeValue strong type={passed ? 'success' : 'danger'} value={observedValue || '<Empty Value>'} />
      </S.Entry>
    </S.AssertionCheckRow>
  );
};

export default AssertionCheckRow;
