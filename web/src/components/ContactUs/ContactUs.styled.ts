import styled from 'styled-components';
import {Modal as AntModal, Button, Typography} from 'antd';
import Plushie from 'assets/plushie.svg';

export const Container = styled.div`
  position: absolute;
  right: 12px;
  bottom: 12px;
  cursor: pointer;
`;

export const PlushieImage = styled.img.attrs({
  src: Plushie,
})``;

export const PulseButtonContainer = styled.div`
  position: absolute;
  left: 2px;
  top: -4px;
`;

export const ModalFooter = styled.div`
  display: flex;
  flex-direction: column;
  gap: 12px;
`;

export const FullWidthButton = styled(Button)`
  && {
    && {
      width: 100%;
      margin: 0px;
    }
  }
`;

export const Header = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
`;

export const Title = styled(Typography.Title)``;

export const Message = styled(Typography.Paragraph)`
  && {
    margin: 0;
  }
`;

export const Modal = styled(AntModal).attrs({
  width: 'auto',
})`
  right: 32px;
  top: calc(100vh - 410px);
  position: absolute;
  height: max-content;
  padding: 0;

  .ant-modal-content {
    width: 340px;
  }

  .ant-modal-footer {
    border: none;
    padding-top: 0;
  }

  .ant-modal-body {
    padding-bottom: 19px;
  }
`;
