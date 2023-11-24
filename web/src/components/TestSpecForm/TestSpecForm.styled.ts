import {DeleteOutlined, ReadOutlined} from '@ant-design/icons';
import {Button, Select as AntSelect, Typography} from 'antd';
import styled from 'styled-components';

export const Select = styled(AntSelect)`
  min-width: 88px;

  > .ant-select-selector {
    min-height: 100%;
  }
`;

export const Check = styled.div`
  display: grid;
  grid-template-columns: repeat(3, 1fr) 14px;
  gap: 8px 12px;
`;

export const AssertionsContainer = styled.div`
  margin-bottom: 24px;
`;

export const Container = styled.div`
  display: flex;
  justify-content: space-between;
  padding-bottom: 8px;
`;

export const FieldsContainer = styled.div`
  display: grid;
  width: 95%;
  grid-template-columns: 1fr auto 1fr;
  gap: 8px;
`;

export const ActionContainer = styled.div`
  display: flex;
  justify-content: center;
  align-self: flex-start;
  flex-basis: 5%;
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
  margin-top: 9px;
`;

export const AssertionForm = styled.div`
  background-color: ${({theme}) => theme.color.white};
  height: 100%;
  overflow-y: auto;
  padding: 24px;
  position: relative;
`;

export const AssertionFormTitle = styled(Typography.Title).attrs({level: 2})``;

export const AssertionFormHeader = styled.div`
  align-items: center;
  display: flex;
  justify-content: space-between;
  margin-bottom: 16px;
`;

export const AssertionFromActions = styled.div`
  display: flex;
  justify-content: flex-end;
  gap: 8px;
`;

export const FormSection = styled.div`
  margin-bottom: 8px;
`;

export const FormSectionRow = styled.div`
  margin-bottom: 8px;
`;

export const FormSectionRow1 = styled.div`
  align-items: center;
  display: flex;
  gap: 12px;
`;

export const FormSectionHeaderSelector = styled.div`
  align-items: center;
  display: flex;
  justify-content: space-between;
`;

export const FormSectionTitle = styled(Typography.Title).attrs({level: 3})<{$noMargin?: boolean}>`
  && {
    margin-bottom: ${({$noMargin}) => ($noMargin ? '0' : '4px')};
  }
`;

export const FormSectionText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
`;

export const SuggestionsContainer = styled.div`
  margin-bottom: 24px;
`;

export const SuggestionsButton = styled(Button)`
  padding: 4px 8px;
`;

export const ReadIcon = styled(ReadOutlined)`
  margin-top: 4px;
`;
