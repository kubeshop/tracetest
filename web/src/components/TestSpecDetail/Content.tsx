import {useEffect, useMemo} from 'react';

import {SemanticGroupNames} from 'constants/SemanticGroupNames.constants';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import AssertionService from 'services/Assertion.service';
import {TAssertionResultEntry} from 'models/AssertionResults.model';
import {useTest} from 'providers/Test/Test.provider';
import useScrollTo from 'hooks/useScrollTo';
import Assertion from './Assertion';
import Header from './Header';
import SpanHeader from './SpanHeader';
import * as S from './TestSpecDetail.styled';

interface IProps {
  onClose(): void;
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry, name: string): void;
  onRevert(originalSelector: string): void;
  onSelectSpan(spanId: string): void;
  selectedSpan?: string;
  testSpec: TAssertionResultEntry;
}

const Content = ({
  onClose,
  onDelete,
  onEdit,
  onRevert,
  onSelectSpan,
  selectedSpan,
  testSpec,
  testSpec: {resultList, selector, spanIds},
}: IProps) => {
  const {
    run: {trace, id: runId},
  } = useTestRun();
  const {
    test: {id: testId},
  } = useTest();
  const {
    isDeleted = false,
    isDraft = false,
    originalSelector = '',
    name = '',
  } = useAppSelector(state => TestSpecsSelectors.selectSpecBySelector(state, selector)) || {};
  const totalPassedChecks = useMemo(() => AssertionService.getTotalPassedChecks(resultList), [resultList]);
  const results = useMemo(() => AssertionService.getResultsHashedBySpanId(resultList), [resultList]);
  const scrollTo = useScrollTo();

  useEffect(() => {
    scrollTo(`assertion-result-${selectedSpan}`);
  }, [scrollTo, selectedSpan]);

  return (
    <>
      <Header
        affectedSpans={spanIds?.length ?? 0}
        assertionsFailed={totalPassedChecks?.false ?? 0}
        assertionsPassed={totalPassedChecks?.true ?? 0}
        isDeleted={isDeleted}
        isDraft={isDraft}
        onClose={onClose}
        onDelete={() => {
          onDelete(testSpec.selector);
          onClose();
        }}
        onEdit={() => {
          onEdit(testSpec, name);
        }}
        onRevert={() => {
          onRevert(originalSelector);
        }}
        selector={selector}
        title={!selector && !name ? 'All Spans' : name}
      />

      {/* // TODO: consider a virtualized list here */}
      {Object.entries(results).map(([spanId, checkResults]) => {
        const span = trace?.spans.find(({id}) => id === spanId);

        return (
          <S.CardContainer
            key={`${testSpec?.id}-${spanId}`}
            size="small"
            title={<SpanHeader onSelectSpan={onSelectSpan} span={span} />}
            type="inner"
            $isSelected={spanId === selectedSpan}
            $type={span?.type ?? SemanticGroupNames.General}
            id={`assertion-result-${spanId}`}
          >
            <S.AssertionsContainer onClick={() => onSelectSpan(span?.id ?? '')}>
              {checkResults.map(checkResult => (
                <Assertion
                  testId={testId}
                  runId={runId}
                  selector={selector}
                  check={checkResult}
                  key={`${checkResult.result.spanId}-${checkResult.assertion}`}
                />
              ))}
            </S.AssertionsContainer>
          </S.CardContainer>
        );
      })}
    </>
  );
};

export default Content;
