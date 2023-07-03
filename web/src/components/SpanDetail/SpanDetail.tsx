import {noop} from 'lodash';
import {useCallback, useMemo, useState} from 'react';

import SearchInput from 'components/SearchInput';
import {
  OtelReference,
  useGetOTELSemanticConventionAttributesInfo,
} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import {CompareOperatorSymbolMap} from 'constants/Operator.constants';
import useSpanData from 'hooks/useSpanData';
import Span from 'models/Span.model';
import TestOutput from 'models/TestOutput.model';
import TestRunOutput from 'models/TestRunOutput.model';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import AssertionService from 'services/Assertion.service';
import SpanService from 'services/Span.service';
import SpanAttributeService from 'services/SpanAttribute.service';
import {TSpanFlatAttribute} from 'types/Span.types';
import {TAnalyzerError, TTestSpecSummary} from 'types/TestRun.types';
import Attributes from './Attributes';
import Header from './Header';
import * as S from './SpanDetail.styled';

export interface IPropsAttributeRow {
  attribute: TSpanFlatAttribute;
  searchText?: string;
  semanticConventions: OtelReference;
  analyzerErrors?: TAnalyzerError[];
  testSpecs?: TTestSpecSummary;
  testOutputs?: TestRunOutput[];
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  onCreateOutput(attribute: TSpanFlatAttribute): void;
}

export interface IPropsSubHeader {
  analyzerErrors?: TAnalyzerError[];
  testSpecs?: TTestSpecSummary;
  testOutputs?: TestRunOutput[];
}

interface IProps {
  onCreateTestSpec?(): void;
  searchText?: string;
  span?: Span;
  AttributeRowComponent: React.ComponentType<IPropsAttributeRow>;
  SubHeaderComponent: React.ComponentType<IPropsSubHeader>;
}

const SpanDetail = ({onCreateTestSpec = noop, searchText, span, AttributeRowComponent, SubHeaderComponent}: IProps) => {
  const {analyzerErrors, testSpecs, testOutputs} = useSpanData(span?.id ?? '');
  const {open} = useTestSpecForm();
  const {onNavigateAndOpen} = useTestOutput();
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
      <Header span={span} />
      <SubHeaderComponent analyzerErrors={analyzerErrors} testSpecs={testSpecs} testOutputs={testOutputs} />
      <S.HeaderDivider />

      <S.SearchContainer data-cy="attributes-search-container">
        <SearchInput placeholder="Search attributes" onSearch={handleOnSearch} width="100%" />
      </S.SearchContainer>

      <Attributes
        attributeList={filteredAttributes}
        searchText={searchText}
        semanticConventions={semanticConventions}
        analyzerErrors={analyzerErrors}
        testSpecs={testSpecs}
        testOutputs={testOutputs}
        onCreateTestSpec={handleCreateTestSpec}
        onCreateOutput={handleCreateOutput}
        AttributeRowComponent={AttributeRowComponent}
      />
    </>
  );
};

export default SpanDetail;
