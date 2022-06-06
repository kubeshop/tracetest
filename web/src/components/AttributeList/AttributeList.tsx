import {useState} from 'react';

import AttributeRow from 'components/AttributeRow';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import * as S from './AttributeList.styled';
import EmptyAttributeList from './EmptyAttributeList';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';

interface IProps {
  assertions?: TResultAssertions;
  attributeList: TSpanFlatAttribute[];
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
}

const AttributeList: React.FC<IProps> = ({assertions, attributeList, onCreateAssertion}) => {
  const [isCopied, setIsCopied] = useState(false);

  const onCopy = (value: string) => {
    TraceAnalyticsService.onAttributeCopy();
    navigator.clipboard.writeText(value);
    setIsCopied(true);
  };

  return attributeList.length ? (
    <S.AttributeList data-cy="attribute-list">
      {attributeList.map(attribute => (
        <AttributeRow
          assertionsFailed={assertions?.[attribute.key]?.failed}
          assertionsPassed={assertions?.[attribute.key]?.passed}
          attribute={attribute}
          isCopied={isCopied}
          key={attribute.key}
          onCopy={onCopy}
          onCreateAssertion={onCreateAssertion}
          setIsCopied={setIsCopied}
        />
      ))}
    </S.AttributeList>
  ) : (
    <EmptyAttributeList />
  );
};

export default AttributeList;
