import {DeleteOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import styled from 'styled-components';

export const Row = styled.div`
  display: flex;
`;

export const Label = styled(Typography.Text).attrs({as: 'div'})`
  margin-bottom: 8px;
`;

export const URLInputContainer = styled.div`
  display: grid;
  align-items: flex-start;
  grid-template-columns: 20% 80%;

  .ant-form-item {
    margin: 0;
  }
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: grid;
  justify-content: center;
  grid-template-columns: 40% 40% 19%;
  margin-bottom: 8px;
`;

export const DeleteIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.md};
`;
