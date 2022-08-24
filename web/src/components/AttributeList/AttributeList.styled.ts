import {QuestionCircleOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const AttributeList = styled.div`
  display: flex;
  flex-direction: column;
`;

export const EmptyAttributeList = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  min-height: 300px;
`;

export const EmptyIcon = styled(QuestionCircleOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: 28px;
  margin-bottom: 16px;
`;

export const EmptyTitle = styled(Typography.Title).attrs({
  level: 3,
  type: 'secondary',
})``;

export const EmptyText = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  max-width: 320px;
  text-align: center;
`;
