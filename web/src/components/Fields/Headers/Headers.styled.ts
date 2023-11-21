import {DeleteOutlined} from '@ant-design/icons';
import styled from 'styled-components';

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
