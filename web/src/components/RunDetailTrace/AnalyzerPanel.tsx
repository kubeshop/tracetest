import TestRun, {isRunStateFinished} from 'models/TestRun.model';
import AnalyzerResult from '../AnalyzerResult/AnalyzerResult';
import SkeletonResponse from '../RunDetailTriggerResponse/SkeletonResponse';
import {RightPanel, PanelContainer} from '../ResizablePanels';

interface IProps {
  run: TestRun;
}

const panel = {
  name: 'ANALYZER',
  maxSize: window.innerWidth / 2 || 650,
  minSize: 25,
  isDefaultOpen: true,
};

const AnalyzerPanel = ({run}: IProps) => (
  <RightPanel panel={panel}>
    {size => (
      <PanelContainer $isOpen={size.isOpen}>
        {isRunStateFinished(run.state) ? <AnalyzerResult result={run.linter} /> : <SkeletonResponse />}
      </PanelContainer>
    )}
  </RightPanel>
);

export default AnalyzerPanel;
