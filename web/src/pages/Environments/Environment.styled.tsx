import {CheckCircleOutlined} from '@ant-design/icons';
import {Button, Typography} from 'antd';
import styled from 'styled-components';

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
    margin: 0;
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

export const ListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 24px;
`;

export const HeaderContainer = styled.div`
  display: flex;
  justify-content: space-between;
  padding-bottom: 8px;
`;

export const VariablesText = styled(Typography)`
  flex-basis: 50%;
`;

export const HeaderText = styled(Typography.Text)`
  flex-basis: 50%;
  font-weight: 600;
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

export const ResultListContainer = styled.div`
  margin: 0 68px 54px 70px;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const NameContainer = styled(TextContainer)`
  display: flex;
  align-items: center;
  gap: 5px;
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const MainHeaderContainer = styled.div`
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  padding: 23px 0;
  width: 100%;
`;

export const InfoIcon = styled(CheckCircleOutlined)`
  color: ${({theme}) => theme.color.text};
  cursor: pointer;
  margin: 4px;
`;
