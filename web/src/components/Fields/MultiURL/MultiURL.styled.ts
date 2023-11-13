import {DeleteOutlined} from '@ant-design/icons';
import styled from 'styled-components';

export const BrokerURLInputContainer = styled.div`
  display: grid;
  align-items: flex-start;
  grid-template-columns: 90% 10%;
  margin-bottom: 8px;

  .ant-form-item {
    margin: 0;
  }
`;

export const DeleteIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.md};
`;
