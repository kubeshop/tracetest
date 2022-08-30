import {DeleteOutlined} from '@ant-design/icons';
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

export const CheckContainer = styled.div`
  display: grid;
  gap: 8px 12px;
  grid-template-columns: repeat(3, 1fr) 14px;
  margin-bottom: 8px;
`;

export const AssertionsContainer = styled.div`
  margin-bottom: 24px;
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
  padding: 24px;
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

export const FormSectionTitle = styled(Typography.Title).attrs({level: 3})<{$noMargin?: boolean}>`
  && {
    margin-bottom: ${({$noMargin}) => ($noMargin ? '0' : '4px')};
  }
`;

export const FormSectionText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
`;
