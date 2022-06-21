import {Button, Typography} from 'antd';
import styled from 'styled-components';

import noResultsIcon from 'assets/SpanAssertionsEmptyState.svg';

export const SpanHeader = styled.div`
  width: 100%;
  align-items: center;
  margin-bottom: 24px;
`;

export const SpanHeaderTitle = styled(Typography.Title).attrs({
  level: 5,
})`
  && {
    margin-bottom: 5px;
  }
`;

export const DetailsContainer = styled.div`
  padding: 24px;
  border: 1px solid rgba(0, 0, 0, 0.06);
  margin-bottom: 16px;
`;

export const DetailsEmptyStateContainer = styled.div`
  display: flex;
  justify-content: center;
  align-items: center;
  margin-top: 40px;
  flex-direction: column;
  gap: 14px;
  overflow-y: auto;
`;

export const DetailsTableEmptyStateIcon = styled.img.attrs({
  src: noResultsIcon,
})``;

export const SpanDetail = styled.div`
  padding: 24px;
  display: flex;
  flex-direction: column;
`;

export const AssertionActionsContainer = styled.div`
  margin-bottom: 24px;
  margin-top: 16px;
  display: flex;
  justify-content: flex-end;
`;

export const AddAssertionButton = styled(Button).attrs({
  type: 'primary',
})`
  && {
    font-weight: 600;
  }
`;

export const Dot = styled.div`
  height: 10px;
  width: 10px;
  margin-left: 5px;
  background-color: #61175e29;
  border-radius: 50%;
  display: inline-block;
`;
