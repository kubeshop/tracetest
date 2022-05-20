import {Drawer} from 'antd';
import {useState, useEffect} from 'react';
import {TTestRun} from 'types/TestRun.types';
import TestResults from '../TestResults';
import * as S from './TraceDrawer.styled';
import TraceDrawerHeader from './TraceDrawerHeader';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import AssertionForm from '../AssertionForm';

interface IProps {
  visiblePortion: number;
  run: TTestRun;
  testId: string;
  onSelectSpan: (spanId: string) => void;
}

const TraceDrawer: React.FC<IProps> = ({run: {id: runId}, run, testId, visiblePortion, onSelectSpan}) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const {isOpen: isAssertionFormOpen, formProps, onSubmit, close} = useAssertionForm();

  useEffect(() => {
    if (isAssertionFormOpen) setIsCollapsed(true);
  }, [isAssertionFormOpen]);

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
        {isAssertionFormOpen ? (
          <AssertionForm
            runId={runId}
            onSubmit={onSubmit}
            testId={testId}
            {...formProps}
            onCancel={() => {
              setIsCollapsed(false);
              close();
            }}
          />
        ) : (
          <TestResults testId={testId} run={run} onSelectSpan={onSelectSpan} />
        )}
      </S.Content>
    </Drawer>
  );
};

export default TraceDrawer;
