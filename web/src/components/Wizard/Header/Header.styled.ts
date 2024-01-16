import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  align-items: center;
  border-bottom: 1px solid ${({theme}) => theme.color.border};
  display: flex;
  justify-content: space-between;
  padding: 16px;
`;

export const ProgressContainer = styled.div`
  width: 300px;
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 8px;
  }
`;

export const Text = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.md};
`;
