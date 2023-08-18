import {Button, Dropdown, Row, Space, Typography} from 'antd';
import styled from 'styled-components';

export const CreateTestButton = styled(Button)`
  font-weight: 600;
`;

export const CreateDropdownButton = styled(Dropdown)``;

export const ActionsContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  margin: 14px 0;
  width: 100%;
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

export const TestListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

export const HeaderContainer = styled.div`
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  padding: 23px 0;
  width: 100%;
`;

export const LoadingContainer = styled(Space)`
  margin-bottom: 24px;
  width: 100%;
`;

export const FiltersContainer = styled.div`
  display: flex;
  gap: 8px;
  align-items: center;
`;

export const ConfigContainer = styled(Row)`
  height: 100%;
`;

export const ConfigContent = styled.div`
  text-align: center;
`;

export const ConfigIcon = styled.img`
  margin-bottom: 25px;
`;

export const ConfigFooter = styled.div`
  margin: 20px 0;
`;
