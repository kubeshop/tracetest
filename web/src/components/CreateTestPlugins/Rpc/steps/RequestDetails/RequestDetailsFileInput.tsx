import {UploadOutlined} from '@ant-design/icons';
import {Upload} from 'antd';
import type {UploadFile} from 'antd/es/upload/interface';
import {noop} from 'lodash';
import * as S from './RequestDetails.styled';

interface IProps {
  onChange?(file?: UploadFile): void;
  value?: UploadFile;
  // accept values https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/file#accept
  accept?: string;
  disabled?: boolean;
  'data-cy'?: string;
}

const RequestDetailsFileInput = ({
  disabled = false,
  accept = '.proto',
  value: file,
  onChange = noop,
  ...props
}: IProps) => (
  <Upload
    disabled={disabled}
    data-cy={props['data-cy']}
    multiple={false}
    fileList={file ? [file] : []}
    onRemove={() => onChange()}
    accept={accept}
    beforeUpload={newFile => {
      onChange(newFile);

      return false;
    }}
  >
    <S.UploadButton data-cy={`${props['data-cy'] || 'upload'}-button`} disabled={disabled} icon={<UploadOutlined />}>
      Choose File
    </S.UploadButton>
  </Upload>
);

export default RequestDetailsFileInput;
