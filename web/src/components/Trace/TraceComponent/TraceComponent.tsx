import * as React from 'react';
import {useState} from 'react';
import {useStoreActions} from 'react-flow-renderer';
import {ISpan} from 'types/Span.types';
import {ITestRunResult} from 'types/TestRunResult.types';
import {ITest} from 'types/Test.types';
import {Diagram, SupportedDiagrams} from 'components/Diagram/Diagram';
import SpanDetail from 'components/SpanDetail';
import {TimelineDrawer} from 'components/Trace/TraceComponent/TimelineDrawer';
import {useHandleOnSpanSelectedCallback} from 'components/Trace/TraceComponent/useHandleOnSpanSelectedCallback';

interface IProps {
  displayError: boolean;
  minHeight: string;
  testResultDetails: ITestRunResult | undefined;
  test?: ITest;
  visiblePortion: number;
}

export const TraceComponent = ({
  displayError,
  visiblePortion,
  minHeight,
  test,
  testResultDetails,
}: IProps): JSX.Element | null => {
  const [selectedSpan, setSelectedSpan] = useState<ISpan | undefined>();
  const addSelected = useStoreActions(actions => actions.addSelectedElements);
  const onSelectSpan = useHandleOnSpanSelectedCallback(addSelected, testResultDetails, setSelectedSpan);
  return !displayError ? (
    <>
      <div style={{display: 'flex', width: '100%', minHeight, maxHeight: minHeight, height: minHeight}}>
        <div style={{flexBasis: '50%', paddingTop: 10, paddingLeft: 10}}>
          <Diagram
            type={SupportedDiagrams.DAG}
            trace={testResultDetails?.trace!}
            onSelectSpan={onSelectSpan}
            selectedSpan={selectedSpan}
          />
        </div>
        <div style={{flexBasis: '50%', overflowY: 'scroll', paddingTop: 10, paddingRight: 10}}>
          <SpanDetail resultId={testResultDetails?.resultId} testId={test?.testId} span={selectedSpan} />
        </div>
      </div>
      <TimelineDrawer
        visiblePortion={visiblePortion}
        testResultDetails={testResultDetails}
        onSelectSpan={onSelectSpan}
        selectedSpan={selectedSpan}
      />
    </>
  ) : null;
};
