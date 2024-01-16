import {Typography} from 'antd';
import styled from 'styled-components';

export const Container = styled.div`
  padding: 24px;
`;

export const Header = styled.div`
  margin-bottom: 16px;
`;

export const Body = styled.div`
  background-color: ${({theme}) => theme.color.white};
  border: 1px solid ${({theme}) => theme.color.border};
  margin-bottom: 24px;
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 6px;
  }
`;

export const Text = styled(Typography.Text)``;
