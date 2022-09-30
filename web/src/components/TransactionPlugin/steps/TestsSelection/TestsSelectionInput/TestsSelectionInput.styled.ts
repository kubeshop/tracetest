import {HolderOutlined, DeleteOutlined} from '@ant-design/icons';
import {SortableHandle} from 'react-sortable-hoc';
import styled from 'styled-components';

export const TestItemContainer = styled.li`
  padding: 0px 10px;
  display: grid;
  grid-template-columns: auto 1fr auto;
  gap: 14px;
  align-items: center;
  border: 1px solid ${({theme}) => theme.color.border};
  border-radius: 2px;
  background: ${({theme}) => theme.color.white};
  height: 32px;
`;

export const DragHandle = SortableHandle(styled(HolderOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  cursor: grab;
`);

export const DeleteIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
`;

export const ItemListContainer = styled.ul`
  list-style: none;
  display: flex;
  gap: 4px;
  flex-direction: column;
  padding: 0;
  margin: 0;
  margin-bottom: 12px;
`;
