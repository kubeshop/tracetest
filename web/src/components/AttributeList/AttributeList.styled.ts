import {QuestionCircleOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const AttributeList = styled.div`
  display: flex;
  flex-direction: column;
  margin-top: 18px;
`;

export const EmptyAttributeList = styled.div`
  display: flex;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  min-height: 450px;
`;

export const EmptyIcon = styled(QuestionCircleOutlined)`
  font-size: 28px;
  color: #687492;
  margin-bottom: 16px;
`;

export const EmptyTitle = styled(Typography.Title).attrs({
  level: 5,
  type: 'secondary',
})``;

export const EmptyText = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  max-width: 400px;
`;
