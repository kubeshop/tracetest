import {Form} from 'antd';
import {TTestVariablesMap} from 'types/Variables.types';
import TestListVariablesInput from './TestVariablesInput/TestListVariablesInput';

interface IProps {
  variables: TTestVariablesMap;
}

const MissingVariablesModalForm = ({variables}: IProps) => {
  return (
    <Form.Item name="variables">
      <TestListVariablesInput variables={variables} />
    </Form.Item>
  );
};

export default MissingVariablesModalForm;
