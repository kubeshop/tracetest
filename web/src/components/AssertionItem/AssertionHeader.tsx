import OperatorService from 'services/Operator.service';
import {TPseudoSelector, TSpanSelector} from 'types/Assertion.types';
import * as S from './AssertionItem.styled';

interface IProps {
  affectedSpans: number;
  failedChecks: number;
  isAdvancedMode: boolean;
  isAdvancedSelector: boolean;
  passedChecks: number;
  pseudoSelector?: TPseudoSelector;
  selectorList: TSpanSelector[];
  title: string;
}

const AssertionHeader = ({
  affectedSpans,
  failedChecks,
  isAdvancedMode,
  isAdvancedSelector,
  passedChecks,
  pseudoSelector,
  selectorList,
  title,
}: IProps) => (
  <S.Column>
    {isAdvancedMode || isAdvancedSelector ? (
      <div>
        <S.HeaderText>{title}</S.HeaderText>
      </div>
    ) : (
      <S.SelectorContainer>
        {selectorList.map(({key, value, operator}) => (
          <S.Selector key={`${key} ${operator} ${value}`}>
            <S.HeaderTextSecondary>
              {key} • {OperatorService.getNameFromSymbol(operator)}
            </S.HeaderTextSecondary>
            <S.HeaderText>{value}</S.HeaderText>
          </S.Selector>
        ))}
        {pseudoSelector && (
          <S.Selector key="pseudo-selector">
            <S.HeaderTextSecondary>pseudo selector</S.HeaderTextSecondary>
            <S.HeaderText>
              {pseudoSelector.selector} {pseudoSelector.number ? `(${pseudoSelector.number})` : ''}
            </S.HeaderText>
          </S.Selector>
        )}
      </S.SelectorContainer>
    )}
    <div>
      {Boolean(passedChecks) && (
        <S.HeaderDetail>
          <S.HeaderDot $passed />
          {passedChecks}
        </S.HeaderDetail>
      )}
      {Boolean(failedChecks) && (
        <S.HeaderDetail>
          <S.HeaderDot $passed={false} />
          {failedChecks}
        </S.HeaderDetail>
      )}
      <S.HeaderDetail>
        <S.HeaderSpansIcon />
        {`${affectedSpans} ${affectedSpans > 1 ? 'spans' : 'span'}`}
      </S.HeaderDetail>
    </div>
  </S.Column>
);

export default AssertionHeader;
