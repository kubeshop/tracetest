import {Progress, Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  padding: 24px;
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 30px;
  }
`;

export const Score = styled(Typography.Title)`
  && {
    font-size: 24px;
    margin-bottom: 0;
  }
`;

export const ScoreContainer = styled.div`
  margin-bottom: 24px;
  text-align: center;
`;

export const RuleContainer = styled.div`
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  margin-bottom: 24px;
  margin-left: 24px;
`;

export const Subtitle = styled(Typography.Title)`
  && {
    margin-bottom: 8px;
  }
`;

export const ScoreProgress = styled(Progress)`
  .ant-progress-inner {
    height: 50px !important;
    width: 50px !important;
  }
`;

export const Column = styled.div`
  display: flex;
  flex-direction: column;
  margin-bottom: 8px;
`;
