import {DeleteOutlined, EditOutlined} from '@ant-design/icons';
import {Tag, Typography} from 'antd';
import styled, {css} from 'styled-components';

export const AssertionCard = styled.div`
  border-radius: 2px;
  border: 1px solid rgba(3, 24, 73, 0.1);
`;

export const Header = styled.div`
  display: flex;
  background: #fbfbff;
  border-bottom: 1px solid rgba(3, 24, 73, 0.1);
  padding: 8px 14px;
  justify-content: space-between;
  border-radius: 2px 2px 0 0;
`;

export const Body = styled.div`
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 9px;
`;

export const SpanCountText = styled(Typography.Text)`
  font-size: 12px;
  margin-right: 14px;
`;

const baseIcon = css`
  font-size: 18px;
  color: #61175e;
  cursor: pointer;
`;

export const EditIcon = styled(EditOutlined)`
  ${baseIcon}
`;

export const DeleteIcon = styled(DeleteOutlined)`
  ${baseIcon};
  margin-left: 12px;
`;

export const ActionsContainer = styled.div`
  display: flex;
  align-items: center;
`;

export const StatusTag = styled(Tag).attrs({
  color: '#61175E',
})``;

export const Selector = styled.div`
  display: flex;
  flex-direction: column;
`;

export const SelectorList = styled.div`
  display: flex;
  gap: 14px;
`;

export const SelectorValueText = styled(Typography.Text)``;

export const SelectorAttributeText = styled(Typography.Text).attrs({
  type: 'secondary',
})`
  font-size: 10px;
  margin-bottom: -3px;
`;
