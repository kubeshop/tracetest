import {useState} from 'react';
import {ITestRunResult} from 'types/TestRunResult.types';
import ResizableDrawer from '../ResizableDrawer/ResizableDrawer';
import TestResults from '../TestResults';

interface IProps {
  visiblePortion: number;
  testResultDetails?: ITestRunResult;
  onSelectSpan: (spanId: string) => void;
}

const TraceDrawer: React.FC<IProps> = ({visiblePortion, onSelectSpan, testResultDetails}) => {
  const [max, setMax] = useState(600);

  return (
    <ResizableDrawer open min={visiblePortion} max={max}>
      <TestResults
        max={max}
        result={testResultDetails!}
        setMax={setMax}
        visiblePortion={visiblePortion}
        onSpanSelected={onSelectSpan}
      />
    </ResizableDrawer>
  );
};

export default TraceDrawer;
