import {Select} from 'antd';
import React from 'react';
import {RequestDefinitionExtended} from './hooks/getRequestsFromCollection';

interface IProps {
  requests: RequestDefinitionExtended[];
  onChange?: (value: string) => void;
}

export const SelectTestFromCollection = ({onChange, requests}: IProps) => {
  return (
    <Select<string>
      style={{width: 490}}
      data-cy="collection-test-select"
      placeholder="Select test from uploaded collection"
      onChange={onChange}
    >
      {requests.map(({id, name}) => (
        <Select.Option data-cy={`postman-test-${id}`} key={id} value={id}>
          {name}
        </Select.Option>
      ))}
    </Select>
  );
};
