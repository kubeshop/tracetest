import {Form, Input} from 'antd';
import Header from 'components/CreateTest/Header';
import RunDetailTriggerResponseFactory from 'components/RunDetailTriggerResponse/RunDetailTriggerResponseFactory';
import RunEvents from 'components/RunEvents';
import FormFactory from 'components/TestPlugins/FormFactory';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import {TriggerTypes} from 'constants/Test.constants';
import {TestState} from 'constants/TestRun.constants';
import {TestRunStage} from 'constants/TestRunEvents.constants';
import useValidateTestDraft from 'hooks/useValidateTestDraft';
import Test from 'models/Test.model';
import TestRun, {isRunStateFinished} from 'models/TestRun.model';
import TestRunEvent from 'models/TestRunEvent.model';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useMemo, useState} from 'react';
import TestService from 'services/Test.service';
import {TDraftTest} from 'types/Test.types';
import * as S from './RunDetailTrigger.styled';
import {useShortcutWithDefault} from '../TestPlugins/hooks/useShortcut';

export const FORM_ID = 'create-test';

interface IProps {
  test: Test;
  run: TestRun;
  runEvents: TestRunEvent[];
  isError: boolean;
}

const RunDetailTrigger = ({test, run: {id, state, triggerResult, triggerTime}, runEvents, isError}: IProps) => {
  const shouldDisplayError = isError || state === TestState.TRIGGER_FAILED;

  const [form] = Form.useForm<TDraftTest>();
  const {isEditLoading: isLoading, onEdit} = useTest();

  const plugin = TriggerTypeToPlugin[test.trigger.type];
  const [isValid, setIsValid] = useState(true);
  const onValidateTest = useValidateTestDraft({pluginName: plugin.name, setIsValid});

  const {run} = useTestRun();
  const stateIsFinished = isRunStateFinished(run.state);

  const initialValues = useMemo(() => TestService.getInitialValues(test), [test]);
  const isDisabled = isLoading || !stateIsFinished;

  useShortcutWithDefault(form);

  return (
    <S.Container>
      <Form<TDraftTest>
        autoComplete="off"
        data-cy="edit-test"
        form={form}
        initialValues={initialValues}
        layout="vertical"
        name={FORM_ID}
        onFinish={values => !isDisabled && onEdit(values)}
        onValuesChange={onValidateTest}
      >
        <Form.Item name="name" hidden>
          <Input type="hidden" value={test.name} />
        </Form.Item>
        <Header isLoading={isDisabled} isValid={isValid} triggerType={test.trigger.type} />

        <S.Body>
          <S.SectionLeft>
            <FormFactory type={test.trigger.type} />
          </S.SectionLeft>

          <S.SectionRight>
            {shouldDisplayError ? (
              <RunEvents events={runEvents} stage={TestRunStage.Trigger} state={state} />
            ) : (
              <RunDetailTriggerResponseFactory
                runId={id}
                state={state}
                testId={test.id}
                triggerResult={triggerResult}
                triggerTime={triggerTime}
                type={triggerResult?.type ?? TriggerTypes.http}
              />
            )}
          </S.SectionRight>
        </S.Body>
      </Form>
    </S.Container>
  );
};

export default RunDetailTrigger;
