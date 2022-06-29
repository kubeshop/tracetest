import {Form, FormInstance, Input} from 'antd';
import {Steps} from 'components/GuidedTour/homeStepList';
import GuidedTourService, {GuidedTours} from 'services/GuidedTour.service';
import {IDemoTestExample} from 'constants/Test.constants';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {SupportedPlugins} from 'constants/Plugins.constants';
import {IBasicDetailsValues} from './BasicDetails';
import BasicDetailsDemoHelper from './BasicDetailsDemoHelper';
import * as S from './BasicDetails.styled';
import useValidate from './hooks/useValidate';

export const FORM_ID = 'create-test';

interface IProps {
  form: FormInstance<IBasicDetailsValues>;
  onSelectDemo(demo: IDemoTestExample): void;
  onSubmit(values: IBasicDetailsValues): void;
  onValidation(isValid: boolean): void;
  selectedDemo?: IDemoTestExample;
}

const BasicDetailsForm = ({form, onSubmit, onSelectDemo, onValidation, selectedDemo}: IProps) => {
  const handleOnValuesChange = useValidate(onValidation);
  const {pluginName} = useCreateTest();

  return (
    <Form
      autoComplete="off"
      data-cy="create-test-modal"
      form={form}
      layout="vertical"
      name={FORM_ID}
      onFinish={onSubmit}
      onValuesChange={handleOnValuesChange}
    >
      <S.GlobalStyle />
      <S.InputContainer>
        {pluginName === SupportedPlugins.REST && (
          <BasicDetailsDemoHelper
            selectedDemo={selectedDemo}
            form={form}
            onSelectDemo={onSelectDemo}
            onValidation={onValidation}
          />
        )}
        <Form.Item
          className="input-name"
          data-cy="create-test-name-input"
          data-tour={GuidedTourService.getStep(GuidedTours.Home, Steps.Name)}
          label="Name"
          name="name"
          rules={[{required: true, message: 'Please enter a test name'}]}
          style={{marginBottom: 0}}
        >
          <Input placeholder="Enter test name" />
        </Form.Item>
        <Form.Item
          className="input-description"
          data-cy="create-test-description-input"
          label="Description"
          name="description"
          style={{marginBottom: 0}}
          rules={[{required: true, message: 'Please enter a test description'}]}
        >
          <Input.TextArea placeholder="Enter a brief description" />
        </Form.Item>
      </S.InputContainer>
    </Form>
  );
};

export default BasicDetailsForm;
