import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  height: 100%;
  padding: 24px;
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin-bottom: 16px;
  }
`;

export const FileName = styled.div`
  margin-bottom: 14px;
`;
