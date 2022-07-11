import {Form} from 'antd';
import React, {Dispatch, SetStateAction} from 'react';
import RequestDetailsFileInput from '../../../../Rpc/steps/RequestDetails/RequestDetailsFileInput';
import {State, useUploadCollectionCallback} from '../hooks/useUploadCollectionCallback';

interface IProps {
  setState: Dispatch<SetStateAction<State>>;
}

export const CollectionFileField = ({setState}: IProps): React.ReactElement => (
  <Form.Item
    rules={[{required: true, message: 'Please enter a request url'}]}
    name="collectionFile"
    label="Upload Postman Collection"
  >
    <RequestDetailsFileInput data-cy="collectionFile" accept=".json" onChange={useUploadCollectionCallback(setState)} />
  </Form.Item>
);
