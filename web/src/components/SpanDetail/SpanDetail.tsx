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
import {TSpan, TSpanFlatAttribute} from 'types/Span.types';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import TestOutput from 'models/TestOutput.model';
import Attributes from './Attributes';
import Header from './Header';
import * as S from './SpanDetail.styled';

interface IProps {
  onCreateTestSpec?(): void;
  searchText?: string;
  span?: TSpan;
}

const SpanDetail = ({onCreateTestSpec = noop, searchText, span}: IProps) => {
  const {open} = useTestSpecForm();
  const {onNavigateAndOpenModal} = useTestOutput();
  const spansResult = useAppSelector(TestSpecsSelectors.selectSpansResult);
  const assertions = useAppSelector(state => TestSpecsSelectors.selectAssertionResultsBySpan(state, span?.id || ''));
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
        selector: {query: selector},
        name: key,
        value: `attr:${key}`,
      });

      onNavigateAndOpenModal(output);
    },
    [onNavigateAndOpenModal, span]
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
      <Header
        span={span}
        totalFailedChecks={span?.id ? spansResult[span.id]?.failed : 0}
        totalPassedChecks={span?.id ? spansResult[span?.id]?.passed : 0}
      />
      <S.HeaderDivider />

      <S.SearchContainer>
        <SearchInput placeholder="Search attributes" onSearch={handleOnSearch} width="100%" />
      </S.SearchContainer>

      <Attributes
        assertions={assertions}
        attributeList={filteredAttributes}
        onCreateTestSpec={handleCreateTestSpec}
        onCreateOutput={handleCreateOutput}
        searchText={searchText}
        semanticConventions={semanticConventions}
      />
    </>
  );
};

export default SpanDetail;
