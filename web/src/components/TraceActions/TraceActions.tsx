import {Button} from 'antd';
import {useTestSpecs} from '../../providers/TestSpecs/TestSpecs.provider';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';
import * as S from './TraceActions.styled';

const TraceActions: React.FC = () => {
  const {specs, publish, cancel} = useTestSpecs();
  const pendingCount = specs.filter(({isDraft}) => isDraft).length;

  return (
    <S.TraceActions>
      <S.ChangesTag data-cy={`trace-actions-pending-count-${pendingCount}`}>
        {pendingCount} pending change(s)
      </S.ChangesTag>
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
    </S.TraceActions>
  );
};

export default TraceActions;
