import {DeleteOutlined, EditOutlined, UndoOutlined} from '@ant-design/icons';
import {Tag, Typography} from 'antd';
import styled, {css} from 'styled-components';

export const AssertionCard = styled.div<{$isSelected: boolean}>`
  border-radius: 2px;
  border: ${({$isSelected, theme}) =>
    $isSelected ? `1px solid ${theme.color.text}` : `1px solid ${theme.color.borderLight}`};
`;

export const Header = styled.div`
  cursor: pointer;
  display: flex;
  background: ${({theme}) => theme.color.background};
  border-bottom: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  padding: 8px 14px;
  justify-content: space-between;
  border-radius: 2px 2px 0 0;
`;

export const Body = styled.div`
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 9px;
  background-color: ${({theme}) => theme.color.white};
`;

export const SpanCountText = styled(Typography.Text)`
  font-size: ${({theme}) => theme.size.sm};
  margin-right: 14px;
`;

const baseIcon = css`
  color: ${({theme}) => theme.color.primary};
  cursor: pointer;
  font-size: ${({theme}) => theme.size.xl};
`;

export const EditIcon = styled(EditOutlined)`
  ${baseIcon};
  margin-left: 12px;
`;

export const DeleteIcon = styled(DeleteOutlined)`
  ${baseIcon};
  margin-left: 12px;
`;

export const UndoIcon = styled(UndoOutlined)`
  ${baseIcon};
`;

export const ActionsContainer = styled.div`
  display: flex;
  align-items: center;
`;

export const StatusTag = styled(Tag)``;

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
  font-size: ${({theme}) => theme.size.xs};
  margin-bottom: -3px;
`;
