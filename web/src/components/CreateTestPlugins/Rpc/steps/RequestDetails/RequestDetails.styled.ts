import {DeleteOutlined} from '@ant-design/icons';
import {Button} from 'antd';
import styled from 'styled-components';

export const InputContainer = styled.div`
  display: grid;
  gap: 26px;
  grid-template-columns: 50%;
  justify-content: start;
`;

export const DoubleInputContainer = styled(InputContainer)`
  display: grid;
  gap: 26px;
  grid-template-columns: 48% 48%;
`;

export const Row = styled.div`
  display: flex;
`;

export const URLInputContainer = styled.div`
  display: flex;
  align-items: end;
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: flex;
  margin-bottom: 8px;
`;

export const UploadButton = styled(Button).attrs({
  type: 'primary',
  ghost: true,
})`
  width: 490px;
`;

export const DeleteIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.md};
`;
