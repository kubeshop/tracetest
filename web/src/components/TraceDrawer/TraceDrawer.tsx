import {Drawer} from 'antd';
import {useEffect} from 'react';
import {TTestRun} from 'types/TestRun.types';
import {useTestDefinition} from '../../providers/TestDefinition/TestDefinition.provider';
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
}

const TraceDrawer: React.FC<IProps> = ({run: {id: runId}, run, testId, visiblePortion, onSelectSpan}) => {
  const {isOpen: isAssertionFormOpen, formProps, onSubmit, close, isCollapsed, setIsCollapsed} = useAssertionForm();
  const {isLoading} = useTestDefinition();

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
        onClick={() => !isAssertionFormOpen && setIsCollapsed(!isCollapsed)}
        run={run}
        isDisabled={isAssertionFormOpen}
        visiblePortion={visiblePortion}
      />
      <S.Content>
        {isLoading && (
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
        {!isLoading && !isAssertionFormOpen && <TestResults testId={testId} run={run} onSelectSpan={onSelectSpan} />}
      </S.Content>
    </Drawer>
  );
};

export default TraceDrawer;
