import RunTopPanel from 'components/RunTopPanel';
import {TTestRun} from 'types/TestRun.types';

interface IProps {
  run: TTestRun;
}

const Run = ({run}: IProps) => {
  return <RunTopPanel run={run} />;
};

export default Run;
