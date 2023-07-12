import {Progress, Typography} from 'antd';
import styled, {DefaultTheme} from 'styled-components';

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

export const Score = styled(Typography.Title)<{$fontSize?: number}>`
  && {
    font-size: ${({$fontSize}) => $fontSize || 12}px;
    margin-bottom: 0;
  }
`;

export const PercentageScoreWrapper = styled.div`
  position: relative;
  display: flex;
  gap: 8px;
  align-items: center;
`;

export const PercentageScore = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.xxxl};
    margin-bottom: 0;
  }
`;

export const ScoreContainer = styled.div`
  margin-bottom: 24px;
  text-align: center;
`;

const getScoreColor = (score: number, theme: DefaultTheme) => {
  if (score >= 90) {
    return theme.color.success;
  }

  if (score >= 60) {
    return theme.color.warningYellow;
  }

  return theme.color.error;
};

interface IScoreProgressProps {
  $height?: string;
  $width?: string;
  $score: number;
}

export const ScoreProgress = styled(Progress).attrs<IScoreProgressProps>(({$score, theme}) => ({
  strokeColor: getScoreColor($score, theme),
}))<IScoreProgressProps>`
  .ant-progress-inner {
    height: ${({$height = '50px'}) => $height} !important;
    width: ${({$width = '50px'}) => $width} !important;
  }

  .ant-progress-circle-trail,
  .ant-progress-circle-path {
    stroke-width: 20px;
  }
`;
