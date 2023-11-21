import {Button, Upload} from 'antd';
import styled from 'styled-components';

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