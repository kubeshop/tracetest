import TestRun, {isRunStateFinished} from 'models/TestRun.model';
import AnalyzerResult from '../AnalyzerResult/AnalyzerResult';
import SkeletonResponse from '../RunDetailTriggerResponse/SkeletonResponse';
import {RightPanel} from '../ResizablePanels';

interface IProps {
  run: TestRun;
}

const panel = {
  isDefaultOpen: true,
};

const AnalyzerPanel = ({run}: IProps) => (
  <RightPanel panel={panel}>
    {isRunStateFinished(run.state) ? <AnalyzerResult result={run.linter} /> : <SkeletonResponse />}
  </RightPanel>
);

export default AnalyzerPanel;
