import {ISpanFlatAttribute} from '../../types/Span.types';
import AttributeRow from '../AttributeRow';
import * as S from './AttributeList.styled';

interface IAttributeListProps {
  attributeList: ISpanFlatAttribute[];
  onCreateAssertion(attribute: ISpanFlatAttribute): void;
}

const AttributeList: React.FC<IAttributeListProps> = ({attributeList, onCreateAssertion}) => {
  return (
    <S.AttributeList>
      {attributeList.map(attribute => (
        <AttributeRow attribute={attribute} key={attribute.key} onCreateAssertion={onCreateAssertion} />
      ))}
    </S.AttributeList>
  );
};

export default AttributeList;
