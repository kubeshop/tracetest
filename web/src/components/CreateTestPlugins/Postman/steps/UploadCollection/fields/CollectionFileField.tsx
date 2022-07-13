import {Form, FormInstance} from 'antd';
import React from 'react';
import RequestDetailsFileInput from '../../../../Rpc/steps/RequestDetails/RequestDetailsFileInput';
import {useUploadCollectionCallback} from '../hooks/useUploadCollectionCallback';
import {IRequestDetailsValues} from '../UploadCollection';

interface IProps {
  form: FormInstance<IRequestDetailsValues>;
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
