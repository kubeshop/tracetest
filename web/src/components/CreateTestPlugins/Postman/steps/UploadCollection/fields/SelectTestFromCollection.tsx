import {Form, FormInstance, Select} from 'antd';
import {useWatch} from 'antd/es/form/Form';
import React, {Dispatch, SetStateAction} from 'react';
import {RequestDefinitionExtended} from 'services/PostmanService.service';
import {useSelectTestCallback} from '../hooks/useSelectTestCallback';
import {IUploadCollectionValues} from '../UploadCollection';

interface IProps {
  form: FormInstance<IUploadCollectionValues>;
  setTransientUrl: Dispatch<SetStateAction<string>>;
}

export const SelectTestFromCollection = ({form, setTransientUrl}: IProps) => {
  const requests = useWatch<RequestDefinitionExtended[]>('requests');
  const variables = useWatch<any[]>('variables');
  return (
    <Form.Item
      rules={[{required: true, message: 'Please enter a request url'}]}
      name="collectionTest"
      label="Select test from Postman Collection"
    >
      <Select<string>
        style={{width: 490}}
        data-cy="collectionTest-select"
        placeholder="Select test from uploaded collection"
        onChange={useSelectTestCallback(form, setTransientUrl, requests, variables)}
      >
        {(requests || []).map(({id, name}, index) => (
          <Select.Option data-cy={`collectionTest-${index}`} key={id} value={id}>
            {name}
          </Select.Option>
        ))}
      </Select>
    </Form.Item>
  );
};
