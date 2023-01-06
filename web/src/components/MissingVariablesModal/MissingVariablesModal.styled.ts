import styled from 'styled-components';
import {Form, Modal as AntModal, Typography} from 'antd';
import {InfoCircleOutlined} from '@ant-design/icons';

export const Modal = styled(AntModal)`
  .ant-modal-body {
    background: ${({theme}) => theme.color.background};
    max-height: 560px;
    overflow-y: scroll;
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
  grid-template-columns: 1fr 1fr;
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

export const Description = styled(Typography.Text).attrs({as: 'p'})`
  && {
    margin: 0;
    margin-bottom: 16px;
  }
`;

export const ToolTipTitle = styled(Typography.Text)`
  && {
    margin: 0;
    font-weight: 600;
  }
`;

export const TestsContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 16px;
`;

export const TestContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
`;

export const FromItem = styled(Form.Item)<{$hasValue: boolean}>`
  margin: 0;

  input {
    border: 1px solid ${({theme, $hasValue}) => !$hasValue && theme.color.error};
  }
`;

export const InfoIcon = styled(InfoCircleOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
  margin: 4px;
`;

export const DetailContainer = styled.div`
  overflow-wrap: break-word;
  width: 270px;
`;

export const TestList = styled.ul`
  list-style-type: none;
  margin: 0;
  padding-left: 0;
`;
