import {Typography} from 'antd';
import styled from 'styled-components';

export const ResultCardList = styled.div``;

export const Header = styled.div`
  display: grid;
  align-items: center;
  grid-template-columns: 300px 300px 100px 220px 40px 40px 40px 1fr;
  gap: 16px;
  padding: 16px 12px;
`;

export const TextContainer = styled.div`
  text-overflow: ellipsis;
  white-space: nowrap;
  overflow: hidden;
`;

export const Text = styled(Typography.Text).attrs({
  strong: true,
})`
  overflow-x: ellipsis;
  font-size: 12px;
`;

export const FailedContainer = styled.div`
  display: flex;
  align-items: center;
  gap: 4px;
`;

export const List = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;
