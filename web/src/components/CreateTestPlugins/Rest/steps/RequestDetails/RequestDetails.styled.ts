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
  display: flex;
  align-items: flex-start;

  .ant-form-item {
    margin: 0;
  }
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: flex;
  margin-bottom: 8px;
`;

export const DeleteIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.md};
`;
