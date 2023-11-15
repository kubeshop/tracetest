import {Form} from 'antd';
import {Overlay} from 'components/Inputs';
import {TDraftTest} from 'types/Test.types';

interface IProps {
  name: string;
  onSubmit(draft: TDraftTest): void;
}

const HeaderForm = ({name, onSubmit}: IProps) => {
  return (
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
        <Overlay />
      </Form.Item>
    </Form>
  );
};

export default HeaderForm;
