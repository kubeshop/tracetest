import {CloseCircleFilled} from '@ant-design/icons';
import styled from 'styled-components';

export const FailedTrace = styled.div`
  height: calc(100vh - 200px);
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;

export const Container = styled.div`
  max-width: 520px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
`;

export const FailedIcon = styled(CloseCircleFilled)`
  color: ${({theme}) => theme.color.error};
  font-size: 32px;
  margin-bottom: 26px;
`;

export const TextContainer = styled(Container)`
  margin-bottom: 26px;
`;
