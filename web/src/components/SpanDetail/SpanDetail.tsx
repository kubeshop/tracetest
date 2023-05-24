import {noop} from 'lodash';
import {useCallback, useMemo, useState} from 'react';

import SearchInput from 'components/SearchInput';
import {useGetOTELSemanticConventionAttributesInfo} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {CompareOperatorSymbolMap} from 'constants/Operator.constants';
import {useAppSelector} from 'redux/hooks';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import AssertionService from 'services/Assertion.service';
import SpanService from 'services/Span.service';
import SpanAttributeService from 'services/SpanAttribute.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {selectOutputsBySpanId} from 'redux/testOutputs/selectors';
import Span from 'models/Span.model';
import TestOutput from 'models/TestOutput.model';
import Attributes from './Attributes';
import Header from './Header';
import * as S from './SpanDetail.styled';

interface IProps {
  onCreateTestSpec?(): void;
  searchText?: string;
  span?: Span;
}

const SpanDetail = ({onCreateTestSpec = noop, searchText, span}: IProps) => {
  const {open} = useTestSpecForm();
  const {onNavigateAndOpen} = useTestOutput();
  const assertions = useAppSelector(state => TestSpecsSelectors.selectAssertionResultsBySpan(state, span?.id || ''));
  const outputs = useAppSelector(state => selectOutputsBySpanId(state, span?.id || ''));
  const [search, setSearch] = useState('');
  const semanticConventions = useGetOTELSemanticConventionAttributesInfo();

  const handleCreateTestSpec = useCallback(
    ({value, key}: TSpanFlatAttribute) => {
      TraceAnalyticsService.onAddAssertionButtonClick();
      const selector = SpanService.getSelectorInformation(span!);

      open({
        isEditing: false,
        selector,
        defaultValues: {
          assertions: [
            {
              left: `attr:${key}`,
              comparator: CompareOperatorSymbolMap.EQUALS,
              right: AssertionService.extractExpectedString(value) || '',
            },
          ],
          selector,
        },
      });

      onCreateTestSpec();
    },
    [onCreateTestSpec, open, span]
  );

  const handleCreateOutput = useCallback(
    ({key}: TSpanFlatAttribute) => {
      TraceAnalyticsService.onAddAssertionButtonClick();
      const selector = SpanService.getSelectorInformation(span!);

      const output = TestOutput({
        selector,
        selectorParsed: {query: selector},
        name: key,
        value: `attr:${key}`,
      });

      onNavigateAndOpen({...output, spanId: span!.id});
    },
    [onNavigateAndOpen, span]
  );

  const handleOnSearch = useCallback((value: string) => {
    setSearch(value);
  }, []);

  const filteredAttributes = useMemo(
    () => SpanAttributeService.filterAttributes(span?.attributeList ?? [], search, semanticConventions),
    [span?.attributeList, search, semanticConventions]
  );

  return (
    <>
      <Header span={span} assertions={assertions} />
      <S.HeaderDivider />

      <S.SearchContainer data-cy="attributes-search-container">
        <SearchInput placeholder="Search attributes" onSearch={handleOnSearch} width="100%" />
      </S.SearchContainer>

      <Attributes
        assertions={assertions}
        attributeList={filteredAttributes}
        onCreateTestSpec={handleCreateTestSpec}
        onCreateOutput={handleCreateOutput}
        searchText={searchText}
        semanticConventions={semanticConventions}
        outputs={outputs}
      />
    </>
  );
};

export default SpanDetail;
