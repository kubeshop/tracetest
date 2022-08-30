import AttributeRow from 'components/AttributeRow';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import * as S from './AttributeList.styled';
import EmptyAttributeList from './EmptyAttributeList';

interface IProps {
  assertions?: TResultAssertions;
  attributeList: TSpanFlatAttribute[];
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  searchText?: string;
}

const AttributeList = ({assertions, attributeList, onCreateTestSpec, searchText}: IProps) => {
  const onCopy = (value: string) => {
    TraceAnalyticsService.onAttributeCopy();
    navigator.clipboard.writeText(value);
  };

  return attributeList.length ? (
    <S.AttributeList data-cy="attribute-list">
      {attributeList.map(attribute => (
        <AttributeRow
          searchText={searchText}
          assertionsFailed={assertions?.[attribute.key]?.failed}
          assertionsPassed={assertions?.[attribute.key]?.passed}
          attribute={attribute}
          key={attribute.key}
          onCopy={onCopy}
          onCreateTestSpec={onCreateTestSpec}
        />
      ))}
    </S.AttributeList>
  ) : (
    <EmptyAttributeList />
  );
};

export default AttributeList;
