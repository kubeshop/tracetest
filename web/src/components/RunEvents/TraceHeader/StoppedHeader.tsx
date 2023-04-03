import {Typography} from 'antd';
import * as S from '../RunEvents.styled';

const StoppedHeader = () => (
  <>
    <S.InfoIcon $isLarge />
    <Typography.Title level={2} type="secondary">
      Test Run Stopped
    </Typography.Title>
    <S.Paragraph type="secondary">The test run was stopped by the user and no trace was detected.</S.Paragraph>
  </>
);

export default StoppedHeader;
