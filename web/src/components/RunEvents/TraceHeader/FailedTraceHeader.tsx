import {Typography} from 'antd';
import * as S from '../RunEvents.styled';

const FailedTraceHeader = () => (
  <>
    <S.ErrorIcon />
    <Typography.Title level={2} type="secondary">
      Trace Fetch Failed
    </Typography.Title>
    <S.Paragraph type="secondary">
      We prepared the breakdown of diagnostic steps and tips to help identify the source of the issue:
    </S.Paragraph>
  </>
);

export default FailedTraceHeader;
