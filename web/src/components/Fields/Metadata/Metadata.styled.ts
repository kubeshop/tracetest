import styled from 'styled-components';
import {DeleteOutlined} from '@ant-design/icons';

export const HeaderContainer = styled.div`
  align-items: center;
  display: flex;
  margin-bottom: 8px;
`;

export const DeleteIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.md};
`;

export const Item = styled.div`
  flex: 1;
  overflow: hidden;
`;
