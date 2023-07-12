import {DeleteOutlined} from '@ant-design/icons';
import {Button, Typography} from 'antd';
import styled from 'styled-components';

export const Title = styled(Typography.Title).attrs({
  level: 3,
})`
  && {
    font-size: ${({theme}) => theme.size.md};
    font-weight: 600;
    margin-bottom: 16px;
  }
`;

export const Subtitle = styled(Typography.Paragraph)`
  && {
    margin-bottom: 8px;
  }
`;

export const TitleContainer = styled.div`
  display: flex;
  justify-content: space-between;
  align-items: center;
`;

export const Container = styled.div`
  margin: 16px 0;
`;

export const SwitchContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
  margin-bottom: 12px;
`;

export const SwitchLabel = styled.label`
  cursor: pointer;
`;

export const ControlsContainer = styled.div`
  margin-top: 16px;
`;

export const OptionsContainer = styled.div``;

export const EntryContainer = styled.div``;

export const ValuesContainer = styled.div`
  align-items: center;
  display: grid;
  justify-content: center;
  gap: 20px;
  grid-template-columns: 1fr 1fr 100px;
  margin-bottom: 4px;
`;

export const DeleteIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.primary};
  font-size: ${({theme}) => theme.size.lg};
`;

export const AddVariableButton = styled(Button)`
  font-weight: 600;
  height: auto;
  padding: 0;
`;

export const DeleteVariableButton = styled(Button)``;

export const VariablesContainer = styled.div`
  margin: 16px 0;
`;

export const TryItButton = styled(Button)`
  && {
    padding: 0 8px;
    background: ${({theme}) => theme.color.white};
    font-weight: 600;

    &:hover,
    &:focus,
    &:active {
      background: ${({theme}) => theme.color.white};
    }
  }
`;
