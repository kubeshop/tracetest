import {useCallback, useState} from 'react';
import {Modal, Typography, FormInstance} from 'antd';

import CreateAssertionForm, {TValues} from './CreateAssertionForm';
import {ISpan} from '../../types/Span.types';
import {IAssertion, IItemSelector} from '../../types/Assertion.types';
import {useAppSelector} from '../../redux/hooks';
import AssertionSelectors from '../../selectors/Assertion.selectors';

interface IProps {
  open: boolean;
  onClose: () => void;
  span: ISpan;
  testId: string;
  resultId: string;
  assertion?: IAssertion;
}

const effectedSpanMessage = (spanCount: number) => {
  if (spanCount <= 1) {
    return `Affects ${spanCount} span`;
  }

  return `Affects ${spanCount} spans`;
};

const CreateAssertionModal = ({testId, span, resultId, open, onClose, assertion}: IProps) => {
  const [form, setForm] = useState<FormInstance<TValues>>();
  const [selectorList, setSelectorList] = useState<IItemSelector[]>([]);

  const onForm = useCallback((formInstance: FormInstance) => {
    setForm(formInstance);
  }, []);

  const onSelectorList = useCallback((selectorListData: IItemSelector[]) => {
    setSelectorList(selectorListData);
  }, []);

  const handleClose = useCallback(() => {
    onClose();
  }, [onClose]);

  const affectedSpanCount = useAppSelector(AssertionSelectors.selectAffectedSpanCount(testId, resultId, selectorList));

  return (
    <Modal
      style={{minHeight: 500}}
      visible={span && open}
      onCancel={handleClose}
      destroyOnClose
      title={
        <div style={{display: 'flex', justifyContent: 'space-between', marginRight: 36}}>
          <Typography.Title level={5}>{assertion ? 'Edit Assertion' : 'Create New Assertion'}</Typography.Title>
          <Typography.Text>{effectedSpanMessage(affectedSpanCount)}</Typography.Text>
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
      />
    </Modal>
  );
};

export default CreateAssertionModal;
