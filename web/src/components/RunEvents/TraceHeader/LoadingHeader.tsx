import {Typography} from 'antd';
import {TRACE_DOCUMENTATION_URL} from 'constants/Common.constants';
import * as S from '../RunEvents.styled';

const LoadingHeader = () => (
  <>
    <S.LoadingIcon />
    <Typography.Title level={3} type="secondary">
      We are working to gather the trace associated with this test run
    </Typography.Title>
    <S.Paragraph type="secondary">
      Want to know more about traces? Head to the official{' '}
      <S.Link href={TRACE_DOCUMENTATION_URL} target="_blank">
        Open Telemetry Documentation
      </S.Link>
    </S.Paragraph>
  </>
);

export default LoadingHeader;
