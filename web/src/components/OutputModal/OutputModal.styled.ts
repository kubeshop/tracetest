import styled from 'styled-components';
import {Modal as AntModal, Typography} from 'antd';

export const Modal = styled(AntModal)`
  .ant-modal-body {
    background: ${({theme}) => theme.color.background};
  }
`;

export const Title = styled(Typography.Title)`
  && {
    font-size: ${({theme}) => theme.size.lg};
    margin: 0;
  }
`;

export const InputContainer = styled.div`
  display: grid;
  gap: 24px;
  grid-template-columns: 100%;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
`;

export const SelectorTitleContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 12px;
`;

export const SelectorLabel = styled(Typography.Text)`
  margin: 0;
`;
