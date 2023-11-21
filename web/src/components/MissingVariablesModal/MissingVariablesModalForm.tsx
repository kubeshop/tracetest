import {Form} from 'antd';
import {TTestVariablesMap} from 'types/Variables.types';
import {TestVariables} from 'components/Inputs';

interface IProps {
  testVariables: TTestVariablesMap;
}

const MissingVariablesModalForm = ({testVariables}: IProps) => {
  return (
    <Form.Item name="variables">
      <TestVariables testVariables={testVariables} />
    </Form.Item>
  );
};

export default MissingVariablesModalForm;
