import {Button} from 'antd';

import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useAppSelector} from 'redux/hooks';
import {selectTestOutputs} from 'redux/testOutputs/selectors';
import TraceAnalyticsService from 'services/Analytics/TestRunAnalytics.service';
import {singularOrPlural} from 'utils/Common';
import * as S from './TestActions.styled';

interface IProps {
  runId: string;
  testId: string;
}

const TestActions = ({runId, testId}: IProps) => {
  const {specs, publish, cancel: onCancelTestSpecs} = useTestSpecs();
  const {onCancel: onCancelTestOutputs} = useTestOutput();
  const outputs = useAppSelector(state => selectTestOutputs(state, testId, runId));
  const pendingSpecs = specs.filter(({isDraft}) => isDraft).length;
  const pendingOutputs = outputs.filter(({isDraft}) => isDraft).length;
  const pendingCount = pendingSpecs + pendingOutputs;

  return (
    <S.Container>
      <S.PendingTag>
        {pendingCount} pending {singularOrPlural('change', pendingCount)}
      </S.PendingTag>
      <Button
        type="link"
        data-cy="trace-actions-revert-all"
        onClick={() => {
          TraceAnalyticsService.onRevertAllClick();
          onCancelTestSpecs();
          onCancelTestOutputs();
        }}
      >
        Revert All
      </Button>
      <Button
        type="primary"
        data-cy="trace-actions-publish"
        onClick={() => {
          TraceAnalyticsService.onPublishClick();
          publish();
        }}
      >
        Publish
      </Button>
    </S.Container>
  );
};

export default TestActions;
