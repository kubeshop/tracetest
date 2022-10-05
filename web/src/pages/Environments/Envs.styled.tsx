import {Button, Typography} from 'antd';
import styled from 'styled-components';
import noResultsIcon from '../../assets/HomeNoResults.svg';

export const CreateEnvironmentButton = styled(Button)``;

export const PageHeader = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  width: 100%;
  margin: 14px 0;
`;

export const TitleText = styled(Typography.Title).attrs({level: 1})`
  && {
    margin: 14px 0;
  }
`;

export const Wrapper = styled.div`
  padding: 0 24px;
  flex-grow: 1;
`;

export const ActionContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
`;

export const NoResultsContainer = styled.div`
  height: 600px;
  display: flex;
  justify-content: center;
  align-items: center;
  flex-direction: column;
`;

export const NoResultsIcon = styled.img.attrs({
  src: noResultsIcon,
})``;

export const NoResultsTitle = styled(Typography.Title)`
  margin-top: 32px;
`;

export const TestListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 24px;
`;

export const VariablesMainContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;
export const VariablesContainer = styled.div`
  align-items: center;
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  border-radius: 2px;
  background: ${({theme}) => theme.color.background};
  display: flex;
  flex-direction: column;
  gap: 16px;
  padding: 12px;
`;

export const NameText = styled(Typography.Title).attrs({ellipsis: true, level: 3})`
  && {
    margin: 0;
  }
`;

export const EnvironmentCard = styled.div<{$isCollapsed: boolean}>`
  box-shadow: 0 4px 8px rgba(153, 155, 168, 0.1);
  background: ${({theme}) => theme.color.white};
  border-left: ${({$isCollapsed, theme}) => $isCollapsed && `2px solid ${theme.color.primary}`};
  border-radius: 2px;
`;

export const InfoContainer = styled.div`
  cursor: pointer;
  display: grid;
  align-items: center;
  grid-template-columns: 20px 1fr 60px 2fr 220px 100px 20px;
  gap: 24px;
  padding: 16px 24px;
`;

export const EnvironmentDetails = styled.div`
  text-align: right;
  width: 100%;
  margin-top: 8px;
`;

export const EnvironmentDetailsLink = styled(Button).attrs({
  type: 'link',
})`
  color: ${({theme}) => theme.color.primary};
  font-weight: 600;
  padding: 0;
`;
