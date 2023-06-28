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

export const InputContainer = styled.div`
  display: grid;
  grid-template-columns: 45% 64px;
  align-items: center;
  max-width: 80%;
`;

export const EditIcon = styled(EditOutlined)`
  && {
    color: ${({theme}) => theme.color.primary};
  }
`;
