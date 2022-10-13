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

export const ValueJson = styled(Typography.Text)`
  cursor: pointer;

  pre {
    margin: 0;
    background: ${({theme}) => theme.color.white};
    border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
    font-size: ${({theme}) => theme.size.sm};
  }
`;

export const ValueText = styled(Typography.Text)`
  margin: 0;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: space-between;
  gap: 8px;
`;

export const SelectorTitleContainer = styled.div`
  display: flex;
  flex-direction: column;
`;

export const SelectorLabel = styled(Typography.Text)`
  margin: 0;
`;

export const SelectedCount = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.xs};
  color: ${({theme}) => theme.color.textSecondary};
  margin: 0;
`;