import {Typography} from 'antd';
import styled from 'styled-components';

export const Body = styled.div`
  min-height: 100px;
  padding: 13px 16px;
`;

export const Container = styled.div`
  background-color: ${({theme}) => theme.color.white};
  width: 300px;
`;

export const Footer = styled.div`
  align-items: center;
  border-top: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  display: flex;
  justify-content: space-between;
  padding: 8px 1px;
`;

export const Header = styled.div`
  align-items: center;
  background: linear-gradient(180deg, #2f1e61 -11.46%, #8b2c53 134.37%);
  display: flex;
  justify-content: space-between;
  padding: 13px 16px;
`;

export const Title = styled(Typography.Title).attrs({level: 3})`
  && {
    color: ${({theme}) => theme.color.white};
    margin-bottom: 0;
  }
`;

export const TitleText = styled(Typography.Text)`
  && {
    color: ${({theme}) => theme.color.white};
    opacity: 0.6;
  }
`;
