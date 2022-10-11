import {Space} from 'antd';
import styled from 'styled-components';

export const EmptyContainer = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  gap: 14px;
  justify-content: center;
  margin: 100px 0;
`;

export const ListContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 4px;
`;

export const LoadingContainer = styled(Space)`
  margin-bottom: 24px;
  width: 100%;
`;
