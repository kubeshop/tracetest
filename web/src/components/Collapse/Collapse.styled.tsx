import {DownOutlined, UpOutlined} from '@ant-design/icons';
import {Collapse} from 'antd';
import styled from 'styled-components';

export const StyledCollapse = styled(Collapse)`
  background-color: ${({theme}) => theme.color.white};
  border: 0;
`;

export const CollapsePanel = styled(Collapse.Panel)`
  background-color: ${({theme}) => theme.color.white};
  border: ${({theme}) => `1px solid ${theme.color.border}`};
  margin-bottom: 12px;

  .ant-collapse-content {
    background-color: ${({theme}) => theme.color.background};
  }
`;

export const CollapseIconContainer = styled.div`
  display: flex;
  position: absolute;
  top: 25%;
  right: 16px;
  border-left: 1px solid ${({theme}) => theme.color.borderLight};
  padding-left: 14px;
  height: 24px;
  align-items: center;
`;

export const DownCollapseIcon = styled(DownOutlined)`
  opacity: 0.5;
  font-size: ${({theme}) => theme.size.xs};
`;

export const UpCollapseIcon = styled(UpOutlined)`
  opacity: 0.5;
  font-size: ${({theme}) => theme.size.xs};
`;
