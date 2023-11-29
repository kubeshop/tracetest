import {InfoCircleOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import Link from 'components/Link';
import * as S from './SkipTraceCollectionInfo.styled';

interface IProps {
  runId: number;
  testId: string;
}

const SkipTraceCollectionInfo = ({runId, testId}: IProps) => {
  return (
    <S.SkipTraceContainer>
      <Typography.Paragraph type="secondary">
        <InfoCircleOutlined /> This test has been set to skip the <b>awaiting trace</b> step. You can change this in{' '}
        <Link to={`/test/${testId}/run/${runId}/trigger?triggerTab=settings`}>
          <Typography.Text type="secondary" underline>
            <b>Settings</b>
          </Typography.Text>
        </Link>
      </Typography.Paragraph>
    </S.SkipTraceContainer>
  );
};

export default SkipTraceCollectionInfo;
