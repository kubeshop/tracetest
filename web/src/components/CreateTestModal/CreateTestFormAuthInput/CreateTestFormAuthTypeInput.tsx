import {Form, FormInstance, Select} from 'antd';
import GuidedTourService, {GuidedTours} from '../../../services/GuidedTour.service';
import {Steps} from '../../GuidedTour/homeStepList';
import {ICreateTestValues} from '../CreateTestForm';
import * as S from '../CreateTestModal.styled';

function methodNaming(method: string | null) {
  switch (method) {
    case 'apiKey':
      return 'API Key';
    case 'bearer':
      return 'Bearer Token';
    case 'basic':
      return 'Basic Auth';
    default:
      return 'No Auth';
  }
}

export const CreateTestFormAuthTypeInput: React.FC<{form: FormInstance<ICreateTestValues>}> = ({form}) => (
  <S.Row>
    <Form.Item
      style={{minWidth: '100%'}}
      initialValue={null}
      label="Authorization Type"
      name={['auth', 'type']}
      valuePropName="type"
    >
      <Select
        className="select-auth-method"
        data-cy="auth-type-select"
        data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Method)}
        dropdownClassName="select-dropdown-auth-type"
        placeholder="No Auth"
        allowClear
        onClear={() => form.resetFields(['auth'])}
        onChange={e => {
          form.resetFields(['auth']);
          form.setFieldsValue({auth: {type: e as any}});
        }}
      >
        {[null, 'apiKey', 'basic', 'bearer'].map(method => (
          <Select.Option data-cy={`auth-type-select-option-${method}`} key={method} value={method}>
            {methodNaming(method)}
          </Select.Option>
        ))}
      </Select>
    </Form.Item>
  </S.Row>
);
