import AttributeRow from 'components/AttributeRow';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import * as S from './AttributeList.styled';

interface IProps {
  assertions?: TResultAssertions;
  attributeList: TSpanFlatAttribute[];
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
}

const AttributeList: React.FC<IProps> = ({assertions, attributeList, onCreateAssertion}) => {
  const onCopy = (value: string) => {
    navigator.clipboard.writeText(value);
  };

  return (
    <S.AttributeList>
      {attributeList.map(attribute => (
        <AttributeRow
          assertionsFailed={assertions?.[attribute.key]?.failed}
          assertionsPassed={assertions?.[attribute.key]?.passed}
          attribute={attribute}
          key={attribute.key}
          onCopy={onCopy}
          onCreateAssertion={onCreateAssertion}
        />
      ))}
    </S.AttributeList>
  );
};

export default AttributeList;
