import {UploadOutlined} from '@ant-design/icons';
import type {UploadFile} from 'antd/es/upload/interface';
import {RcFile} from 'antd/lib/upload';
import {noop} from 'lodash';
import * as S from './FileUpload.styled';

interface IProps {
  onChange?(file?: RcFile): void;
  value?: UploadFile;
  // accept values https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/file#accept
  accept?: string;
  disabled?: boolean;
  'data-cy'?: string;
}

const FileUpload = ({disabled = false, accept = '.proto', value: file, onChange = noop, ...props}: IProps) => (
  <S.UploadContainer
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
  </S.UploadContainer>
);

export default FileUpload;
