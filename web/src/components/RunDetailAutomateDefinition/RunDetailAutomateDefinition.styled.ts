import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  padding: 24px;
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin-bottom: 16px;
    font-weight: 700;
  }
`;

export const SubtitleContainer = styled.div`
  margin-bottom: 8px;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  margin-top: 16px;
`;
