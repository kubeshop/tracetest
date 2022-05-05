import {Button, Typography} from 'antd';
import {DISCORD_URL, GITHUB_ISSUES_URL} from '../../constants/Common.constants';
import * as S from './FailedTrace.styled';

interface FailedTraceProps {
  onReRun(): void;
  onEdit(): void;
}

const FailedTrace: React.FC<FailedTraceProps> = ({onReRun, onEdit}) => {
  return (
    <S.FailedTrace>
      <S.Container>
        <S.FailedIcon />
        <S.TextContainer>
          <Typography.Title level={3}>Test Run Failed</Typography.Title>
          <Typography.Text type="secondary">Information explaining the state the test failed at.</Typography.Text>
          <Typography.Text type="secondary">
            Please let us know about this issue - <a href={GITHUB_ISSUES_URL}>create an issue</a> or contact us via{' '}
            <a href={DISCORD_URL}>Discord</a>.
          </Typography.Text>
          <Typography.Text type="secondary">We will check it out and respond to you.</Typography.Text>
        </S.TextContainer>
        <S.ButtonContainer>
          <Button type="primary" ghost onClick={onEdit}>
            Edit Test
          </Button>
          <Button type="primary" ghost onClick={onReRun}>
            Rerun Test
          </Button>
        </S.ButtonContainer>
      </S.Container>
    </S.FailedTrace>
  );
};

export default FailedTrace;
