import {Form} from 'antd';
import { TTestVariablesMap } from 'types/Variables.types';
import TestListVariablesInput from './TestVariablesInput/TestListVariablesInput';

interface IProps {
  testVariables: TTestVariablesMap;
}

const MissingVariablesModalForm = ({testVariables}: IProps) => {
  return (
    <Form.Item name="variables">
      <TestListVariablesInput testVariables={testVariables} />
    </Form.Item>
  );
};

export default MissingVariablesModalForm;
