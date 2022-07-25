import {DeleteOutlined} from '@ant-design/icons';
import {Button, Upload} from 'antd';
import styled from 'styled-components';

export const Row = styled.div`
  display: flex;
`;

export const HeaderContainer = styled.div`
  align-items: center;
  display: flex;
  margin-bottom: 8px;
`;

export const UploadContainer = styled(Upload)`
  .ant-upload {
    width: 100%;
  }
`;

export const UploadButton = styled(Button).attrs({
  type: 'primary',
  ghost: true,
})`
  width: 100%;
`;

export const DeleteIcon = styled(DeleteOutlined)`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme}) => theme.size.md};
`;
