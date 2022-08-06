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

export const AdvancedSelectorInputContainer = styled.div`
  display: flex;
  //align-items: start;
  //grid-template-columns: 850px 170px;
  //gap: 14px;
  //max-width: 85%;
  //width: 85%;
  //grid-template-columns: 100%;
`;

export const PseudoSelector = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 4px;
`;

export const Check = styled.div`
  display: flex;
  //grid-template-columns: repeat(3, 210px) 1fr;
  //gap: 4px;
  //align-items: start;
  //margin-bottom: 16px;
`;

export const AddCheckButton = styled(Button).attrs({
  type: 'link',
})`
  padding: 0;
  font-weight: 600;
`;

export const DeleteCheckIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: pointer;
`;

export const AssertionForm = styled.div`
  background-color: ${({theme}) => theme.color.white};
  border: ${({theme}) => `1px solid ${theme.color.borderLight}`};
  padding: 16px;
  width: 100%;
`;

export const AssertionFormTitle = styled(Typography.Title).attrs({level: 3})``;

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
  max-width: 85%;
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

export const AffectedSpansLabel = styled(Typography.Text)`
  color: ${({theme}) => theme.color.primary};
`;
