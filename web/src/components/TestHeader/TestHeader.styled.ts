import {LeftOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const TestHeader = styled.div`
  background-color: #fff;
  display: flex;
  gap: 14px;
  align-items: center;
  justify-content: space-between;
  padding: 19px 24px;
  border-bottom: 1px solid rgba(3, 24, 73, 0.1);
  width: 100%;
`;

export const TestName = styled(Typography.Title).attrs({
  level: 4,
})`
  && {
    margin: 0;
    font-weight: 400;
    font-size: 18px;
  }
`;

export const TestUrl = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  && {
    margin: 0;
    align-self: flex-end;
    font-size: 12px;
  }
`;

export const StateText = styled(Typography.Text)`
  && {
    margin-right: 8px;
    color: #8c8c8c;
  }
`;

export const BackIcon = styled(LeftOutlined)`
  cursor: pointer;
  font-size: 18px;
`;

export const Content = styled.div`
  display: flex;
  gap: 14px;
  align-items: center;
`;

export const StateContainer = styled.div`
  align-items: center;
  display: flex;
  justify-self: flex-end;
  cursor: pointer;
`;
