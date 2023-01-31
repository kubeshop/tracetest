import {useCallback} from 'react';

import LoadingSpinner from 'components/LoadingSpinner';
import TestSpecs from 'components/TestSpecs';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useAppSelector} from 'redux/hooks';
import TestSpecsSelectors from 'selectors/TestSpecs.selectors';
import AssertionAnalyticsService from 'services/Analytics/AssertionAnalytics.service';
import {TAssertionResultEntry} from 'models/AssertionResults.model';
import Header from './Header';
import * as S from './TestResults.styled';

interface IProps {
  onDelete(selector: string): void;
  onEdit(assertionResult: TAssertionResultEntry): void;
  onRevert(originalSelector: string): void;
}

const TestResults = ({onDelete, onEdit, onRevert}: IProps) => {
  const {isLoading, assertionResults, setSelectedSpec} = useTestSpecs();
  const {selectedSpan, onSetFocusedSpan, onSelectSpan} = useSpan();
  const {totalFailedSpecs, totalPassedSpecs} = useAppSelector(TestSpecsSelectors.selectTotalSpecs);

  const handleOpen = useCallback(
    (selector: string) => {
      AssertionAnalyticsService.onAssertionClick();
      const testSpec = assertionResults?.resultList?.find(specResult => specResult.selector === selector);

      onSetFocusedSpan('');
      onSelectSpan(testSpec?.spanIds[0] || '');
      setSelectedSpec(testSpec?.selector);
    },
    [assertionResults?.resultList, onSelectSpan, onSetFocusedSpan, setSelectedSpec]
  );

  return (
    <S.Container>
      <Header selectedSpan={selectedSpan!} totalFailedSpecs={totalFailedSpecs} totalPassedSpecs={totalPassedSpecs} />

      {(isLoading || !assertionResults) && (
        <S.LoadingContainer>
          <LoadingSpinner />
        </S.LoadingContainer>
      )}

      {!isLoading && Boolean(assertionResults) && (
        <TestSpecs
          assertionResults={assertionResults!}
          onDelete={onDelete}
          onEdit={onEdit}
          onOpen={handleOpen}
          onRevert={onRevert}
        />
      )}
    </S.Container>
  );
};

export default TestResults;
