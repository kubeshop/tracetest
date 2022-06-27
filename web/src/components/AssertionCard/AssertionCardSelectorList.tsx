import OperatorService from '../../services/Operator.service';
import {TPseudoSelector, TSpanSelector} from '../../types/Assertion.types';
import * as S from './AssertionCard.styled';

interface IProps {
  selectorList: TSpanSelector[];
  pseudoSelector?: TPseudoSelector;
}

const AssertionCardSelectorList: React.FC<IProps> = ({selectorList, pseudoSelector}) => {
  return (
    <S.SelectorList>
      {selectorList.map(({key, value, operator}) => (
        <S.Selector key={`${key} ${operator} ${value}`}>
          <S.SelectorAttributeText>
            {key} • {OperatorService.getNameFromSymbol(operator)}
          </S.SelectorAttributeText>
          <S.SelectorValueText>{value}</S.SelectorValueText>
        </S.Selector>
      ))}
      {pseudoSelector && (
        <S.Selector key="pseudo-selector">
          <S.SelectorAttributeText>pseudo selector</S.SelectorAttributeText>
          <S.SelectorValueText>
            {pseudoSelector.selector} {pseudoSelector.number ? `(${pseudoSelector.number})` : ''}
          </S.SelectorValueText>
        </S.Selector>
      )}
    </S.SelectorList>
  );
};

export default AssertionCardSelectorList;
