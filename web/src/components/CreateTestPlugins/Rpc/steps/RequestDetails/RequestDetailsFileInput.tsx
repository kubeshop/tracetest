import {UploadOutlined} from '@ant-design/icons';
import type {UploadFile} from 'antd/es/upload/interface';
import {noop} from 'lodash';
import {Upload} from 'antd';
import * as S from './RequestDetails.styled';

interface IProps {
  onChange?(file?: UploadFile): void;
  value?: UploadFile;
  // https://developer.mozilla.org/en-US/docs/Web/HTML/Element/input/file#accept
  accept?: string;
}

const RequestDetailsFileInput = ({accept = '.proto', value: file, onChange = noop}: IProps) => (
  <Upload
    multiple={false}
    fileList={file ? [file] : []}
    onRemove={() => onChange()}
    accept={accept}
    beforeUpload={newFile => {
      onChange(newFile);

      return false;
    }}
  >
    <S.UploadButton icon={<UploadOutlined />}>Choose File</S.UploadButton>
  </Upload>
);

export default RequestDetailsFileInput;
