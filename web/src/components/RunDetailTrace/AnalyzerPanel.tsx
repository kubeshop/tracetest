import TestRun, {isRunStateFinished} from 'models/TestRun.model';
import Trace from 'models/Trace.model';
import AnalyzerResult from '../AnalyzerResult/AnalyzerResult';
import SkeletonResponse from '../RunDetailTriggerResponse/SkeletonResponse';
import {RightPanel, PanelContainer} from '../ResizablePanels';

interface IProps {
  run: TestRun;
}

const panel = {
  name: 'ANALYZER',
  maxSize: 650,
  minSize: 25,
};

const AnalyzerPanel = ({run}: IProps) => {
  return (
    <RightPanel panel={panel}>
      {size => (
        <PanelContainer $isOpen={size.isOpen}>
          {isRunStateFinished(run.state) ? (
            <AnalyzerResult result={run.linter} trace={run?.trace ?? Trace({})} />
          ) : (
            <SkeletonResponse />
          )}
        </PanelContainer>
      )}
    </RightPanel>
  );
};

export default AnalyzerPanel;
