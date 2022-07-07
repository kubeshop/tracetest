import {CloseCircleFilled} from '@ant-design/icons';
import styled from 'styled-components';

export const Container = styled.div`
  height: calc(100vh - 200px);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;

export const Content = styled.div`
  display: flex;
  max-width: 800px;
  padding: 24px;
`;

export const Icon = styled(CloseCircleFilled)`
  color: ${({theme}) => theme.color.error};
  font-size: 32px;
`;
