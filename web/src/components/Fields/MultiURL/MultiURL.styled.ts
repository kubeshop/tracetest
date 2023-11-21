import {DeleteOutlined} from '@ant-design/icons';
import styled from 'styled-components';

export const URLContainer = styled.div`
  align-items: start;
  display: flex;

  .ant-form-item {
    margin: 0;
  }
`;

export const DeleteIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.md};
`;

export const AddURLContainer = styled.div`
  position: absolute;
  right: 8px;
  top: 4px;
`;
