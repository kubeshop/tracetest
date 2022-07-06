import {LeftOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const CreateTestHeader = styled.div`
  background-color: ${({theme}) => theme.color.white};
  display: flex;
  gap: 14px;
  align-items: center;
  justify-content: space-between;
  padding: 19px 24px;
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  width: 100%;
`;

export const Name = styled(Typography.Title).attrs({level: 1})`
  && {
    margin: 0;
  }
`;

export const BackIcon = styled(LeftOutlined)`
  cursor: pointer;
  font-size: ${({theme}) => theme.size.xl};
`;

export const Content = styled.div`
  display: flex;
  gap: 14px;
  align-items: center;
`;

export const Row = styled.div`
  display: flex;
`;
