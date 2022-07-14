import AssertionForm from 'components/AssertionForm';
import AssertionList from 'components/AssertionList';
import {useAssertionForm} from 'components/AssertionForm/AssertionForm.provider';
import LoadingSpinner from 'components/LoadingSpinner';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestDefinition} from 'providers/TestDefinition/TestDefinition.provider';
import {TTestRun} from 'types/TestRun.types';
import Header from './Header';
import * as S from './RunBottomPanel.styled';

interface IProps {
  run: TTestRun;
  testId: string;
}

const RunBottomPanel = ({run: {id: runId}, run, testId}: IProps) => {
  const {isOpen: isAssertionFormOpen, formProps, onSubmit, close} = useAssertionForm();
  const {isLoading, assertionResults} = useTestDefinition();
  const {selectedSpan, onSelectSpan} = useSpan();

  return (
    <>
      <Header
        assertionResults={assertionResults}
        isDisabled={isAssertionFormOpen}
        run={run}
        selectedSpan={selectedSpan!}
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
            <AssertionList assertionResults={assertionResults!} onSelectSpan={onSelectSpan} />
          )}
        </S.Content>
      </S.Container>
    </>
  );
};

export default RunBottomPanel;
