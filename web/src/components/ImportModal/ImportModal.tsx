import {Button, Form} from 'antd';
import {useCallback, useState} from 'react';
import {TDraftTest} from 'types/Test.types';
import {ImportTypes} from 'constants/Test.constants';
import ImportService from 'services/Import.service';
import {useDashboard} from 'providers/Dashboard/Dashboard.provider';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import * as S from './ImportModal.styled';
import {ImportSelector} from '../Inputs';
import ImportFactory from '../TestPlugins/ImportFactory';
import Tip from './Tip';

interface IProps {
  isOpen: boolean;
  onClose(): void;
}

const FORM_ID = 'import-test';

const ImportModal = ({isOpen, onClose}: IProps) => {
  const [form] = Form.useForm<TDraftTest>();
  const type = Form.useWatch('importType', form) || ImportTypes.curl;
  const {onUpdateDraftTest, onUpdatePlugin} = useCreateTest();
  const {navigate} = useDashboard();

  const [isValid, setIsValid] = useState(false);
  const handleChange = useCallback(
    async (values: TDraftTest) => {
      const valid = await ImportService.validateDraft(type, values);
      setIsValid(valid);
    },
    [type]
  );

  const handleImport = useCallback(
    async (values: TDraftTest) => {
      const {draft, plugin} = await ImportService.getRequest(type, values);

      onUpdatePlugin(plugin);
      onUpdateDraftTest(draft);
      navigate(`/test/create/${plugin.type}`);
    },
    [navigate, onUpdateDraftTest, onUpdatePlugin, type]
  );

  return (
    <S.Modal
      footer={
        <>
          <Button onClick={onClose}>Cancel</Button>
          <Button type="primary" disabled={!isValid} onClick={() => form.submit()}>
            Import
          </Button>
        </>
      }
      onCancel={onClose}
      title="Import to Tracetest"
      visible={isOpen}
      centered
    >
      <Form<TDraftTest>
        form={form}
        autoComplete="off"
        data-cy="import-test-modal"
        layout="vertical"
        name={FORM_ID}
        onFinish={handleImport}
        initialValues={{importType: ImportTypes.curl}}
        onValuesChange={(_: any, values) => handleChange(values)}
      >
        <S.Container>
          <Form.Item name="importType">
            <ImportSelector />
          </Form.Item>

          <ImportFactory type={type} />
        </S.Container>
        <Tip />
      </Form>
    </S.Modal>
  );
};

export default ImportModal;
