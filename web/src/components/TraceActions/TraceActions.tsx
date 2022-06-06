import {Button} from 'antd';
import {useTestDefinition} from '../../providers/TestDefinition/TestDefinition.provider';
import TraceAnalyticsService from '../../services/Analytics/TraceAnalytics.service';
import * as S from './TraceActions.styled';

const TraceActions: React.FC = () => {
  const {definitionList, publish, cancel} = useTestDefinition();
  const pendingCount = definitionList.filter(({isDraft}) => isDraft).length;

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
          publish;
        }}
      >
        Publish
      </Button>
    </S.TraceActions>
  );
};

export default TraceActions;
