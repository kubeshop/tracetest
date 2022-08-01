import AttributeRow from 'components/AttributeRow';
import {useState} from 'react';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import {useSpan} from 'providers/Span/Span.provider';
import {useGuidedTour} from 'providers/GuidedTour/GuidedTour.provider';
import * as S from './AttributeList.styled';
import EmptyAttributeList from './EmptyAttributeList';

interface IProps {
  assertions?: TResultAssertions;
  attributeList: TSpanFlatAttribute[];
  onCreateAssertion(attribute: TSpanFlatAttribute): void;
}

const AttributeList: React.FC<IProps> = ({assertions, attributeList, onCreateAssertion}) => {
  const [isCopied, setIsCopied] = useState(false);
  const {searchText} = useSpan();

  const onCopy = (value: string) => {
    TraceAnalyticsService.onAttributeCopy();
    navigator.clipboard.writeText(value);
    setIsCopied(true);
  };

  const {
    tour: {isOpen, currentStep},
  } = useGuidedTour();

  const getShouldDisplayActions = (index: number) => isOpen && currentStep === 3 && !index;

  return attributeList.length ? (
    <S.AttributeList data-cy="attribute-list">
      {attributeList.map((attribute, index) => (
        <AttributeRow
          searchText={searchText}
          shouldDisplayActions={getShouldDisplayActions(index)}
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
