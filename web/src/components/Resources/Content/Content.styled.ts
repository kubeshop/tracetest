import {Typography} from 'antd';
import styled from 'styled-components';

export const Body = styled.div`
  display: flex;
  gap: 24px;
  width: 100%;
`;

export const Container = styled.div`
  padding: 0;
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 6px;
  }
`;

export const Text = styled(Typography.Text)``;

export const Link = styled(Typography.Link)`
  font-weight: bold;
`;

export const Card = styled.div`
  background-color: ${({theme}) => theme.color.white};
  border-radius: 2px;
  display: flex;
  flex: 1;
  flex-direction: column;
  gap: 14px;
  padding: 24px;
`;

export const Icon = styled.img`
  width: 64px;
  height: auto;
`;
