import TestRun from 'models/TestRun.model';
import {RightPanel} from '../ResizablePanels';
import Specs from './Specs';

interface IProps {
  run: TestRun;
}

const panel = {
  name: 'SPECS',
  isDefaultOpen: true,
  order: 0,
};

const SpecsPanel = ({run}: IProps) => (
  <RightPanel panel={panel}>
    <Specs run={run} />
  </RightPanel>
);

export default SpecsPanel;
