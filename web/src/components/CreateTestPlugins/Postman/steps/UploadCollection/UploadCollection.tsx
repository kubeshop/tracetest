import {Form} from 'antd';
import {RcFile} from 'antd/lib/upload';
import * as Step from 'components/CreateTestPlugins/Step.styled';
import {useCallback, useEffect} from 'react';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import {useCreateTest} from 'providers/CreateTest/CreateTest.provider';
import {IPostmanValues} from 'types/Test.types';
import {HTTP_METHOD} from 'constants/Common.constants';
import {ComponentNames} from 'constants/Plugins.constants';
import UploadCollectionForm from './UploadCollectionForm';
import {useUploadCollectionCallback} from './hooks/useUploadCollectionCallback';

const UploadCollection = () => {
  const {onNext, pluginName, draftTest, onIsFormValid} = useCreateTest();
  const [form] = Form.useForm<IPostmanValues>();
  const {url = '', body = '', method = HTTP_METHOD.GET, collectionFile, collectionTest} = draftTest as IPostmanValues;

  const onValidate = useValidateTestDraft({pluginName, setIsValid: onIsFormValid});
  const getCollectionValues = useUploadCollectionCallback(form);

  const handleOnSubmit = useCallback(
    (values: IPostmanValues) => {
      onNext(values);
    },
    [onNext]
  );

  const currentUrl = Form.useWatch('url', form);

  const onRefreshData = useCallback(async () => {
    form.setFieldsValue({url, body, method: method as HTTP_METHOD, collectionFile, collectionTest});
    getCollectionValues(collectionFile as RcFile);
    onIsFormValid(true);
  }, [form, url, body, method, collectionFile, collectionTest, getCollectionValues, onIsFormValid]);

  useEffect(() => {
    onRefreshData();
  }, [onRefreshData]);

  const onValidateUrlChange = useCallback(async () => {
    try {
      await form.validateFields();
      onIsFormValid(true);
    } catch (err) {
      onIsFormValid(false);
    }
  }, [form, onIsFormValid]);

  useEffect(() => {
    onValidateUrlChange();
  }, [currentUrl, onValidateUrlChange]);

  return (
    <Step.Step>
      <Step.FormContainer>
        <Step.Title>Method Selection Information</Step.Title>
        <Form<IPostmanValues>
          id={ComponentNames.UploadCollection}
          autoComplete="off"
          form={form}
          layout="vertical"
          initialValues={{url: '', requests: [], variables: []}}
          onFinish={handleOnSubmit}
          onValuesChange={onValidate}
        >
          <UploadCollectionForm form={form} />
        </Form>
      </Step.FormContainer>
    </Step.Step>
  );
};

export default UploadCollection;
