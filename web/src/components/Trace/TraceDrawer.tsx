import {Drawer} from 'antd';
import {useState} from 'react';
import {ITestRunResult} from 'types/TestRunResult.types';
import TestResults from '../TestResults';

interface IProps {
  visiblePortion: number;
  testResultDetails?: ITestRunResult;
  onSelectSpan: (spanId: string) => void;
}

const TraceDrawer: React.FC<IProps> = ({visiblePortion, onSelectSpan, testResultDetails}) => {
  const [isCollapsed, setIsCollapsed] = useState(false);

  return (
    <Drawer
      placement="bottom"
      closable={false}
      mask={false}
      visible
      height={isCollapsed ? '420px' : visiblePortion}
      style={{overflow: 'hidden'}}
      bodyStyle={{overflow: 'hidden', padding: 0}}
    >
      <TestResults
        onHeaderClick={() => setIsCollapsed(!isCollapsed)}
        result={testResultDetails!}
        visiblePortion={visiblePortion}
        onSpanSelected={onSelectSpan}
      />
    </Drawer>
  );
};

export default TraceDrawer;
