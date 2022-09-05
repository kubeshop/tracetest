import {Button} from 'antd';

import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import TraceAnalyticsService from 'services/Analytics/TraceAnalytics.service';
import {singularOrPlural} from 'utils/Common';
import * as S from './TestActions.styled';

const TestActions = () => {
  const {specs, publish, cancel} = useTestSpecs();
  const pendingCount = specs.filter(({isDraft}) => isDraft).length;

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
          cancel();
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
