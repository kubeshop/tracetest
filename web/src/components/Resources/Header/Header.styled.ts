import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  align-items: center;
  display: flex;
  flex: 1;
  justify-content: space-between;
  padding: 4px 0;
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 8px;
  }
`;

export const Text = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.md};
`;
