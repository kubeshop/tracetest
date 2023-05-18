import {Progress, Typography} from 'antd';
import styled from 'styled-components';

export const ScoreWrapper = styled.div`
  position: relative;
`;

export const ScoreTexContainer = styled.div`
  position: absolute;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  width: 100%;
  height: 100%;
`;

export const Score = styled(Typography.Title)`
  && {
    font-size: 12px;
    margin-bottom: 0;
  }
`;

export const ScoreContainer = styled.div`
  margin-bottom: 24px;
  text-align: center;
`;

export const ScoreProgress = styled(Progress)<{$height?: string, $width?: string }>`
  .ant-progress-inner {
    height: ${({$height = "50px"}) => $height} !important;
    width: ${({$width = "50px"}) => $width} !important;
  }

  .ant-progress-circle-trail,
  .ant-progress-circle-path {
    stroke-width: 20px;
  }
`;
