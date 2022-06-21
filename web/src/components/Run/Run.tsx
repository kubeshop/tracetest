import RunBottomPanel from 'components/RunBottomPanel';
import {RunLayoutProvider} from 'components/RunLayout';
import RunTopPanel from 'components/RunTopPanel';
import {TTest} from 'types/Test.types';
import {TTestRun} from 'types/TestRun.types';

interface IProps {
  displayError: boolean;
  run: TTestRun;
  test?: TTest;
}

const Run = ({displayError, run, test}: IProps) => {
  if (displayError) {
    return null;
  }

  return (
    <RunLayoutProvider
      bottomPanel={<RunBottomPanel run={run} testId={test?.id!} />}
      topPanel={<RunTopPanel run={run} />}
    />
  );
};

export default Run;
