import {CloseCircleOutlined, EditOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
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
  border-radius: 2px 2px 0px 0px;
`;

export const Body = styled.div`
  padding: 12px 14px;
  display: flex;
  flex-direction: column;
  gap: 9px;
`;

export const SelectorListText = styled(Typography.Text).attrs({
  strong: true,
})`
  margin-right: 14px;
`;

export const SpanCountText = styled(Typography.Text)`
  font-size: 12;
`;
const baseIcon = css`
  font-size: 18px;
  color: #61175e;
  cursor: pointer;
`;

export const EditIcon = styled(EditOutlined)`
  ${baseIcon}
`;

export const DeleteIcon = styled(CloseCircleOutlined)`
  ${baseIcon}
  margin-left: 12px;
`;
