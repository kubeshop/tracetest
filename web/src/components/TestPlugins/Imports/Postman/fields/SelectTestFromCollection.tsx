import {Form, Select} from 'antd';
import {useWatch} from 'antd/es/form/Form';
import {RequestDefinitionExtended} from 'services/Importers/Postman.service';
import {IPostmanValues, TDraftTestForm} from 'types/Test.types';
import {useSelectTestCallback} from '../hooks/useSelectTestCallback';

interface IProps {
  form: TDraftTestForm<IPostmanValues>;
}

export const SelectTestFromCollection = ({form}: IProps) => {
  const requests = useWatch<RequestDefinitionExtended[]>('requests');
  const variables = useWatch<any[]>('variables');

  return (
    <Form.Item
      rules={[{required: true, message: 'No test selected yet'}]}
      name="collectionTest"
      label="Select test from Postman collection"
    >
      <Select<string>
        data-cy="collectionTest-select"
        placeholder="Select test from uploaded collection"
        onChange={useSelectTestCallback(form, requests, variables)}
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
