import {useState} from 'react';
import {ITestRunResult} from 'types/TestRunResult.types';
import ResizableDrawer from '../ResizableDrawer/ResizableDrawer';
import {TraceTimeline} from './TraceTimeline';
import {ISpan} from '../../types/Span.types';

interface IProps {
  visiblePortion: number;
  testResultDetails: ITestRunResult | undefined;
  onSelectSpan: (spanId: string) => void;
  selectedSpan?: ISpan;
}

export const TimelineDrawer = ({
  visiblePortion,
  onSelectSpan,
  selectedSpan,
  testResultDetails,
}: IProps): JSX.Element => {
  const [max, setMax] = useState(600);
  return (
    <ResizableDrawer open min={visiblePortion} max={max}>
      <TraceTimeline
        max={max}
        setMax={setMax}
        visiblePortion={visiblePortion}
        trace={testResultDetails?.trace!}
        onSelectSpan={onSelectSpan}
        selectedSpan={selectedSpan}
      />
    </ResizableDrawer>
  );
};
