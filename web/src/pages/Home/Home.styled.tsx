import {Button, Typography} from 'antd';
import styled from 'styled-components';
import noResultsIcon from '../../assets/HomeNoResults.svg';

export const CreateTestButton = styled(Button)``;

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

export const NoResultsTitle = styled(Typography.Title).attrs({
  level: 3,
})`
  margin-top: 32px;
`;

export const TestListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 24px;
`;
