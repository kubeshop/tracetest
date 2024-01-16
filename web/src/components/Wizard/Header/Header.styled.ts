import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  padding: 16px;
  border-bottom: 1px solid ${({theme}) => theme.color.border};
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 12px;
  }
`;

export const Text = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.sm};
`;
