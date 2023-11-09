import {EditOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import styled from 'styled-components';

export const Overlay = styled.div`
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
`;

export const SaveButton = styled(Button)``;

export const InputContainer = styled.label`
  & > input {
    min-width: max-content;
    border-top: none;
    border-left: none;
    border-right: none;
    padding: 0px;

    word-wrap: break-word;
    word-break: break-all;
  }
`;

export const EditIcon = styled(EditOutlined)`
  && {
    color: ${({theme}) => theme.color.primary};
  }
`;
