import {CheckCircleOutlined} from '@ant-design/icons';
import {Tabs, Typography} from 'antd';
import styled, {css} from 'styled-components';

const defaultHeight = '100vh - 106px - 60px - 40px';

export const DataStoreListContainer = styled(Tabs)`
  height: calc(${defaultHeight} - 50px);

  && {
    .ant-tabs-content-holder {
      width: 1px;
    }

    .ant-tabs-tab {
      margin: 0 !important;
      padding: 0;
    }

    .ant-tabs-nav-list {
      gap: 16px;
    }
  }
`;

export const DataStoreItemContainer = styled.div<{$isDisabled: boolean; $isSelected: boolean}>`
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 22px;
  cursor: pointer;

  ${({$isDisabled}) =>
    $isDisabled &&
    css`
      cursor: not-allowed;
      opacity: 0.5;
    `}
`;

export const DataStoreName = styled(Typography.Text)<{$isSelected: boolean}>`
  && {
    color: ${({theme, $isSelected}) => ($isSelected ? theme.color.primary : theme.color.text)};
    font-weight: ${({$isSelected}) => ($isSelected ? 700 : 400)};
  }
`;

export const InfoIcon = styled(CheckCircleOutlined)`
  color: ${({theme}) => theme.color.text};
  cursor: pointer;
  margin: 4px;
`;
