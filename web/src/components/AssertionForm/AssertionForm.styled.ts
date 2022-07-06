import {DeleteOutlined} from '@ant-design/icons';
import {Button, Select as AntSelect, Typography} from 'antd';
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
  align-items: start;
  grid-template-columns: 850px 170px;
  gap: 14px;
`;

export const AdvancedSelectorInputContainer = styled(SelectorInputContainer)`
  grid-template-columns: 80%;
`;

export const PseudoSelector = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
`;

export const Check = styled.div`
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
  background-color: #fff;
`;

export const AssertionFormTitle = styled(Typography.Text).attrs({
  as: 'p',
})`
  margin: 0;
  margin-bottom: 16px;
  font-weight: 600;
`;

export const AssertionFormHeader = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const AssertionFromActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
`;

export const CheckActions = styled.div`
  display: flex;
  align-items: center;
  gap: 14px;
  margin-left: 14px;
  height: 100%;
`;

export const AdvancedSelectorContainer = styled.div`
  display: flex;
  gap: 8px;
  align-items: center;
  margin-bottom: 8px;
  max-width: 80%;
`;

export const AffectedSpansContainer = styled.div`
  display: flex;
  gap: 4px;
  align-items: center;
`;

export const ReferenceLink = styled(Typography.Text)`
  margin: 0;
  margin-left: auto;
`;
