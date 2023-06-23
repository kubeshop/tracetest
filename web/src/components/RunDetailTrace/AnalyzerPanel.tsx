import TestRun, {isRunStateFinished} from 'models/TestRun.model';
import Trace from 'models/Trace.model';
import {TPanel, TPanelComponentProps} from '../ResizablePanels/ResizablePanels';
import AnalyzerResult from '../AnalyzerResult/AnalyzerResult';
import SkeletonResponse from '../RunDetailTriggerResponse/SkeletonResponse';
import * as S from './RunDetailTrace.styled';

type TProps = TPanelComponentProps & {
  run: TestRun;
};

const AnalyzerPanel = ({run, size: {isOpen}}: TProps) => {
  return isRunStateFinished(run.state) ? (
    <S.PanelContainer $isOpen={isOpen}>
      <AnalyzerResult result={run.linter} trace={run?.trace ?? Trace({})} />
    </S.PanelContainer>
  ) : (
    <SkeletonResponse />
  );
};

export const getAnalyzerPanel = (run: TestRun, order = 1): TPanel => ({
  name: `ANALYZER_${order}`,
  maxSize: 720,
  minSize: 15,
  isDefaultOpen: true,
  position: 'right',
  component: ({size}) => <AnalyzerPanel size={size} run={run} />,
});

export default AnalyzerPanel;
