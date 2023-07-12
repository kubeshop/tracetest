import {Typography} from 'antd';
import * as S from '../RunEvents.styled';

const FailedTriggerHeader = () => (
  <>
    <S.ErrorIcon />
    <Typography.Title level={2} type="secondary">
      Test Trigger Failed
    </Typography.Title>
    <S.Paragraph type="secondary">
      The test failed in the Trigger stage, review the Trigger tab to see the breakdown of diagnostic steps.
    </S.Paragraph>
  </>
);

export default FailedTriggerHeader;
