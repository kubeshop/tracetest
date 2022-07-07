import {LeftOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const TestHeader = styled.div`
  background-color: ${({theme}) => theme.color.white};
  display: flex;
  gap: 14px;
  align-items: center;
  justify-content: space-between;
  padding: 19px 24px;
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  width: 100%;
`;

export const TestName = styled(Typography.Title).attrs({level: 1})`
  && {
    margin: 0;
  }
`;

export const TestUrl = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  && {
    font-size: ${({theme}) => theme.size.sm};
    margin: 0;
  }
`;

export const StateText = styled(Typography.Text)`
  && {
    margin-right: 8px;
    color: ${({theme}) => theme.color.textSecondary};
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

export const StateContainer = styled.div`
  align-items: center;
  display: flex;
  justify-self: flex-end;
  cursor: pointer;
`;

export const Row = styled.div`
  display: flex;
`;

export const RightSection = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
`;
