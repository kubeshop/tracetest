import {DeleteOutlined} from '@ant-design/icons';
import {Button, Select as AntSelect, Space, Typography} from 'antd';
import styled from 'styled-components';

export const Select = styled(AntSelect)`
  min-width: 88px;
  > .ant-select-selector {
    min-height: 100%;
  }
`;

export const SelectorContainer = styled.div`
  display: flex;
  align-items: center;
`;

export const SelectorInputContainer = styled.div`
  display: grid;
  align-items: center;
  grid-template-columns: 850px 170px;
  gap: 14px;
`;

export const PseudoSelector = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
`;

export const Check = styled(Space)`
  display: grid;
  grid-template-columns: repeat(3, 210px) 1fr;
  gap: 4px;
  align-items: start;
  margin-bottom: 16px;
`;

export const AddCheckButton = styled(Button).attrs({
  type: 'link',
})`
  padding: 0;
  font-weight: 600;
`;

export const DeleteCheckIcon = styled(DeleteOutlined).attrs({
  style: {color: 'rgb(140, 140, 140)'},
})`
  cursor: pointer;
`;

export const AssertionForm = styled.div`
  width: 100%;
  padding: 16px;
  border: 1px solid rgba(3, 24, 73, 0.1);
`;

export const AssertionFormTitle = styled(Typography.Text).attrs({
  as: 'p',
})`
  margin: 0;
  margin-bottom: 16px;
  font-weight: 600;
`;

export const AssertionFromActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
`;
