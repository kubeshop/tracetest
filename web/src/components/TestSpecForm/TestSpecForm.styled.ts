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

export const AssertionsContainer = styled.div`
  margin-bottom: 24px;
`;

export const Container = styled.div`
  display: flex;
  justify-content: space-between;
  padding-bottom: 8px;
`;

export const FieldsContainer = styled.div`
  display: flex;
  justify-content: space-between;
  flex-basis: 95%;
`;

export const ActionContainer = styled.div`
  display: flex;
  justify-content: center;
  align-self: center;
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

export const ExpectedInputContainer = styled.div`
  width: 0;
  flex-basis: 50%;
  padding-left: 8px;

  && {
    .cm-editor {
      overflow: hidden;
      display: flex;
      border-radius: 2px;
      font-size: ${({theme}) => theme.size.md};
      outline: 1px solid grey;
      height: 32px;
      font-family: SFPro, serif;
      outline: 1px solid #CDD1DB;
    }

    .cm-content {
      display: flex;
      align-items: center;
    }
    .cm-scroller {
      overflow: hidden;
    }
    .cm-line {
      padding: 0;
      span {
        font-family: SFPro, serif;
      }
    }
  }
`;
