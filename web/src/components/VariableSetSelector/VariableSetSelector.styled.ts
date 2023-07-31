import {FormOutlined, PlusOutlined} from '@ant-design/icons';
import styled from 'styled-components';

export const VarsSelectorEntryIcon = styled(FormOutlined)`
  && {
    color: ${({theme}) => theme.color.textSecondary};
    width: 13px;
    height: 13px;
  }
`;

export const VarsSelectorEntryContainer = styled.div`
  display: flex;
  align-items: center;
  justify-content: space-between;
  min-width: 180px;
  max-width: 180px;
`;

export const AddVarsContainer = styled(VarsSelectorEntryContainer)`
  justify-content: start;
  gap: 10px;
  color: ${({theme}) => theme.color.primary};
  font-weight: 600;
`;

export const AddVarsIcon = styled(PlusOutlined)`
  && {
    color: ${({theme}) => theme.color.primary};
  }
`;
