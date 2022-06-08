import AssertionForm from 'components/AssertionForm';
import {useAssertionForm} from 'components/AssertionForm/AssertionFormProvider';
import LoadingSpinner from 'components/LoadingSpinner';
import TestResults from 'components/TestResults';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {TSpan} from 'types/Span.types';
import {TTestRun} from 'types/TestRun.types';
import Header from './Header';
import * as S from './RunBottomPanel.styled';

interface IProps {
  onSelectSpan: (spanId: string) => void;
  run: TTestRun;
  selectedSpan: TSpan;
  testId: string;
}

const RunBottomPanel = ({onSelectSpan, run: {id: runId}, run, selectedSpan, testId}: IProps) => {
  const {isOpen: isAssertionFormOpen, formProps, onSubmit, close} = useAssertionForm();
  const {isLoading, assertionResults} = useTestDefinition();

  return (
    <>
      <Header
        assertionResults={assertionResults}
        isDisabled={isAssertionFormOpen}
        run={run}
        selectedSpan={selectedSpan}
      />
      <S.Container id="assertions-container">
        <S.Content>
          {(isLoading || !assertionResults) && (
            <S.LoadingSpinnerContainer>
              <LoadingSpinner />
            </S.LoadingSpinnerContainer>
          )}
          {!isLoading && isAssertionFormOpen && (
            <AssertionForm
              runId={runId}
              onSubmit={onSubmit}
              testId={testId}
              {...formProps}
              onCancel={() => {
                close();
              }}
            />
          )}
          {!isLoading && !isAssertionFormOpen && Boolean(assertionResults) && (
            <TestResults testId={testId} assertionResults={assertionResults!} onSelectSpan={onSelectSpan} />
          )}
        </S.Content>
      </S.Container>
    </>
  );
};

export default RunBottomPanel;
