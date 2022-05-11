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
`;

export const FailedIcon = styled(CloseCircleFilled).attrs({
  style: {
    color: '#e84749',
    fontSize: 32,
  },
})`
  margin-bottom: 26px;
`;

export const TextContainer = styled(Container)`
  margin-bottom: 26px;
`;

export const ButtonContainer = styled.div`
  display: flex;
  gap: 8px;
`;
