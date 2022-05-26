import {Drawer} from 'antd';
import {useEffect} from 'react';
import {TTestRun} from 'types/TestRun.types';
import {useTestDefinition} from '../../providers/TestDefinition/TestDefinition.provider';
import {TSpan} from '../../types/Span.types';
import AssertionForm from '../AssertionForm';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import LoadingSpinner from '../LoadingSpinner';
import TestResults from '../TestResults';
import * as S from './TraceDrawer.styled';
import TraceDrawerHeader from './TraceDrawerHeader';

interface IProps {
  visiblePortion: number;
  run: TTestRun;
  testId: string;
  onSelectSpan: (spanId: string) => void;
  selectedSpan: TSpan;
}

const TraceDrawer: React.FC<IProps> = ({run: {id: runId}, run, testId, visiblePortion, onSelectSpan, selectedSpan}) => {
  const {isOpen: isAssertionFormOpen, formProps, onSubmit, close, setIsCollapsed, isCollapsed} = useAssertionForm();
  const {isLoading, assertionResults} = useTestDefinition();

  useEffect(() => {
    if (isAssertionFormOpen) setIsCollapsed(true);
  }, [isAssertionFormOpen, setIsCollapsed]);

  return (
    <Drawer
      placement="bottom"
      closable={false}
      mask={false}
      visible
      data-cy="trace-drawer"
      height={isCollapsed ? '420px' : visiblePortion}
      style={{overflow: 'hidden'}}
      bodyStyle={{overflow: 'hidden', padding: 0}}
    >
      <TraceDrawerHeader
        onClick={() => setIsCollapsed(!isCollapsed)}
        run={run}
        assertionResults={assertionResults}
        isDisabled={isAssertionFormOpen}
        visiblePortion={visiblePortion}
        selectedSpan={selectedSpan}
      />
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
    </Drawer>
  );
};

export default TraceDrawer;
