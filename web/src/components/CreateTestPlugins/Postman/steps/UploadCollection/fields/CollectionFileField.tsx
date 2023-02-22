import {Form} from 'antd';
import RequestDetailsFileInput from 'components/CreateTestPlugins/Grpc/steps/RequestDetails/RequestDetailsFileInput';
import React from 'react';
import {IPostmanValues, TDraftTestForm} from '../../../../../../types/Test.types';
import {useUploadCollectionCallback} from '../hooks/useUploadCollectionCallback';

interface IProps {
  form: TDraftTestForm<IPostmanValues>;
}

export const CollectionFileField = ({form}: IProps): React.ReactElement => (
  <Form.Item
    rules={[{required: true, message: 'No file selected yet'}]}
    name="collectionFile"
    label="Upload Postman Collection"
  >
    <RequestDetailsFileInput data-cy="collectionFile" accept=".json" onChange={useUploadCollectionCallback(form)} />
  </Form.Item>
);
