import {LeftOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const CreateTestHeader = styled.div`
  background-color: #fff;
  display: flex;
  gap: 14px;
  align-items: center;
  justify-content: space-between;
  padding: 19px 24px;
  border-bottom: 1px solid rgba(3, 24, 73, 0.1);
  width: 100%;
`;

export const Name = styled(Typography.Title).attrs({
  level: 4,
})`
  && {
    margin: 0;
    font-size: 16px;
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

export const Row = styled.div`
  display: flex;
`;
