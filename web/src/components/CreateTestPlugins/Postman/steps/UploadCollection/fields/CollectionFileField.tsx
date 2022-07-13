import {Form, FormInstance} from 'antd';
import RequestDetailsFileInput from 'components/CreateTestPlugins/Rpc/steps/RequestDetails/RequestDetailsFileInput';
import React from 'react';
import {useUploadCollectionCallback} from '../hooks/useUploadCollectionCallback';
import {IUploadCollectionValues} from '../UploadCollection';

interface IProps {
  form: FormInstance<IUploadCollectionValues>;
}

export const CollectionFileField = ({form}: IProps): React.ReactElement => (
  <Form.Item
    rules={[{required: true, message: 'Please enter a request url'}]}
    name="collectionFile"
    label="Upload Postman Collection"
  >
    <RequestDetailsFileInput data-cy="collectionFile" accept=".json" onChange={useUploadCollectionCallback(form)} />
  </Form.Item>
);
