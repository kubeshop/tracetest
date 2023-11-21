import styled from 'styled-components';
import {Modal as AntModal, Typography} from 'antd';

export const Modal = styled(AntModal)`
  top: 50px;

  .ant-modal-body {
    background: ${({theme}) => theme.color.background};
    max-height: calc(100vh - 250px);
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

export const TestContainer = styled.div`
  display: flex;
  flex-direction: column;
  gap: 8px;
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
