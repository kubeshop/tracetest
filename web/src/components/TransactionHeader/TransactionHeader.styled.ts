import {LeftOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import {Link as RRLink} from 'react-router-dom';
import styled from 'styled-components';

export const BackIcon = styled(LeftOutlined)`
  cursor: pointer;
  font-size: ${({theme}) => theme.size.lg};
`;

export const Container = styled.div`
  align-items: center;
  background-color: ${({theme}) => theme.color.white};
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  display: flex;
  justify-content: space-between;
  padding: 6px 24px;
  width: 100%;
`;

export const Section = styled.div`
  align-items: center;
  display: flex;
  flex: 1;
  gap: 14px;
`;

export const SectionRight = styled(Section)`
  justify-content: flex-end;
`;

export const Text = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  && {
    font-size: ${({theme}) => theme.size.sm};
    margin: 0;
  }
`;

export const Title = styled(Typography.Title).attrs({level: 2})`
  && {
    margin: 0;
  }
`;

export const StateContainer = styled.div`
  align-items: center;
  display: flex;
  justify-self: flex-end;
  cursor: pointer;
`;

export const StateText = styled(Typography.Text)`
  && {
    margin-right: 8px;
    color: ${({theme}) => theme.color.textSecondary};
  }
`;

export const LinksContainer = styled.div`
  a:first-child {
    border-top-left-radius: 2px;
    border-bottom-left-radius: 2px;
  }

  a:last-child {
    border-top-right-radius: 2px;
    border-bottom-right-radius: 2px;
  }
`;

export const Link = styled(RRLink)<{$isActive?: boolean}>`
  && {
    background-color: ${({theme, $isActive}) => ($isActive ? theme.color.primary : theme.color.white)};
    color: ${({theme, $isActive}) => ($isActive ? theme.color.white : theme.color.primary)};
    font-weight: 600;
    padding: 7px 16px;

    border: 1px solid rgba(3, 24, 73, 0.1);

    &:hover,
    &:visited,
    &:focused {
      color: ${({theme, $isActive}) => $isActive && theme.color.white};
    }
  }
`;
