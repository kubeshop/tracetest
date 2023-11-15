import {Form} from 'antd';
import {Overlay} from 'components/Inputs';
import {TDraftTest} from 'types/Test.types';

interface IProps {
  name: string;
  onSubmit(draft: TDraftTest): void;
  isDisabled: boolean;
}

const HeaderForm = ({name, onSubmit, isDisabled}: IProps) => (
  <Form<TDraftTest>
    autoComplete="off"
    initialValues={{name}}
    name="edit-test-name"
    onValuesChange={(changedValues: any, draft: TDraftTest) => {
      if (draft.name === name) return;
      onSubmit(draft);
    }}
  >
    <Form.Item name="name" noStyle>
      <Overlay isDisabled={isDisabled} />
    </Form.Item>
  </Form>
);

export default HeaderForm;
