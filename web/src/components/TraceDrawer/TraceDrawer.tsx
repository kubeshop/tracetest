import {Drawer} from 'antd';
import {useState, useEffect} from 'react';
import {useSelector} from 'react-redux';
import {ITestRunResult} from 'types/TestRunResult.types';
import TestResultSelectors from '../../selectors/TestResult.selectors';
import TestResults from '../TestResults';
import * as S from './TraceDrawer.styled';
import TraceDrawerHeader from './TraceDrawerHeader';
import {useAssertionForm} from '../AssertionForm/AssertionFormProvider';
import AssertionForm from '../AssertionForm';

interface IProps {
  visiblePortion: number;
  result: ITestRunResult;
  onSelectSpan: (spanId: string) => void;
}

const TraceDrawer: React.FC<IProps> = ({result: {resultId, testId}, result, visiblePortion, onSelectSpan}) => {
  const [isCollapsed, setIsCollapsed] = useState(false);
  const {isOpen: isAssertionFormOpen, formProps, onSubmit, close} = useAssertionForm();

  const traceResultList = useSelector(TestResultSelectors.selectTestResultList(resultId));

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
        result={result}
        isDisabled={isAssertionFormOpen}
        visiblePortion={visiblePortion}
      />
      <S.Content>
        {isAssertionFormOpen ? (
          <AssertionForm
            resultId={resultId}
            onSubmit={onSubmit}
            testId={testId}
            {...formProps}
            onCancel={() => {
              setIsCollapsed(false);
              close();
            }}
          />
        ) : (
          <TestResults assertionResultList={traceResultList} result={result} onSelectSpan={onSelectSpan} />
        )}
      </S.Content>
    </Drawer>
  );
};

export default TraceDrawer;
