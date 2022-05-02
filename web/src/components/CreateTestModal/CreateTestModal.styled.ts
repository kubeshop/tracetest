import {Typography} from 'antd';
import styled from 'styled-components';

export const Wrapper = styled.div`
  .test-modal {
    overflow: hidden !important;
    display: flex;
  }

  .test-modal .ant-modal {
    top: 50px;
  }

  .test-modal .ant-modal-content {
    height: min(calc(100vh - 16%), 620px);
  }

  .test-modal .ant-modal-body {
    position: absolute;
    min-height: fit-content;
    height: calc(100% - 110px);
    width: 100%;
    overflow-y: auto;
  }

  .test-modal .ant-modal-footer {
    position: absolute;
    left: 0;
    right: 0;
    bottom: 0;
  }

  .method-select .ant-select-selector {
    background-color: #fafafa !important;
  }

  .method-select-item .ant-select-item-option-selected {
    background-color: #fafafa !important;
  }
`;

export const DropdownText = styled(Typography.Text).attrs({
  as: 'a',
})``;

export const DemoTextContainer = styled.div`
  margin-bottom: 24px;
`;
