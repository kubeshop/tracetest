import {useCallback, useState} from 'react';
import {Modal, Typography, FormInstance} from 'antd';

import CreateAssertionForm, {TValues} from './CreateAssertionForm';
import {getEffectedSpansCount} from '../../entities/Assertion/Assertion.service';
import { TSpan } from '../../entities/Span/Span.types';
import { TTrace } from '../../entities/Trace/Trace.types';
import { TAssertion, TItemSelector } from '../../entities/Assertion/Assertion.types';

interface IProps {
  open: boolean;
  onClose: () => void;
  span: TSpan;
  testId: string;
  trace: TTrace;
  assertion?: TAssertion;
}

const effectedSpanMessage = (spanCount: number) => {
  if (spanCount <= 1) {
    return `Effects ${spanCount} span`;
  }

  return `Effects ${spanCount} spans`;
};

const CreateAssertionModal = ({testId, span, trace, open, onClose, assertion}: IProps) => {
  const [form, setForm] = useState<FormInstance<TValues>>();
  const [selectorList, setSelectorList] = useState<TItemSelector[]>([]);

  const onForm = useCallback((formInstance: FormInstance) => {
    setForm(formInstance);
  }, []);

  const onSelectorList = useCallback((selectorListData: TItemSelector[]) => {
    setSelectorList(selectorListData);
  }, []);

  const handleClose = useCallback(() => {
    onClose();
  }, [onClose]);

  const effectedSpanCount = getEffectedSpansCount(trace, selectorList);

  return (
    <Modal
      style={{minHeight: 500}}
      visible={span && open}
      onCancel={handleClose}
      destroyOnClose
      title={
        <div style={{display: 'flex', justifyContent: 'space-between', marginRight: 36}}>
          <Typography.Title level={5}>{assertion ? 'Edit Assertion' : 'Create New Assertion'}</Typography.Title>
          <Typography.Text>{effectedSpanMessage(effectedSpanCount)}</Typography.Text>
        </div>
      }
      onOk={form?.submit}
      okButtonProps={{
        type: 'default',
      }}
      okText="Save"
    >
      <CreateAssertionForm
        assertion={assertion}
        onCreate={handleClose}
        onForm={onForm}
        onSelectorList={onSelectorList}
        span={span}
        testId={testId}
        trace={trace}
      />
    </Modal>
  );
};

export default CreateAssertionModal;
