import {TSpanFlatAttribute} from '../../types/Span.types';
import AttributeRow from '../AttributeRow';
import * as S from './AttributeList.styled';
import EmptyAttributeList from './EmptyAttributeList';

interface IAttributeListProps {
  attributeList: TSpanFlatAttribute[];
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
}

const AttributeList: React.FC<IAttributeListProps> = ({attributeList, onCreateAssertion}) => {
  return attributeList.length ? (
    <S.AttributeList data-cy="attribute-list">
      {attributeList.map(attribute => (
        <AttributeRow attribute={attribute} key={attribute.key} onCreateAssertion={onCreateAssertion} />
      ))}
    </S.AttributeList>
  ) : (
    <EmptyAttributeList />
  );
};

export default AttributeList;
