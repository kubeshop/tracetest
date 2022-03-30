import {Button, Typography} from 'antd';
import styled from 'styled-components';
import noResultsIcon from '../../assets/HomeNoResults.svg';

export const CreateTestButton = styled(Button)``;

export const PageHeader = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
  width: 100%;
  margin-bottom: 32px;
`;

export const TitleText = styled(Typography.Title).attrs({
  level: 3,
})``;

export const Wrapper = styled.div`
  padding: 0 24px;
`;

export const NoResultsContainer = styled.div`
  border: 1px solid #f0f0f0;
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
