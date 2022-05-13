import {useCallback, useState} from 'react';
import {Modal, Typography, FormInstance} from 'antd';

import CreateAssertionForm, {TValues} from './CreateAssertionForm';
import {ISpan, ISpanFlatAttribute} from '../../types/Span.types';
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
  defaultAttributeList?: ISpanFlatAttribute[];
}

const affectedSpanMessage = (spanCount: number) => {
  if (spanCount === 1) {
    return `Affect ${spanCount} span`;
  }

  return `Affects ${spanCount} spans`;
};

const CreateAssertionModal = ({testId, span, resultId, open, onClose, assertion, defaultAttributeList = []}: IProps) => {
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

  const affectedSpanList = useAppSelector(AssertionSelectors.selectAffectedSpanList(testId, resultId, selectorList));

  return (
    <Modal
      style={{minHeight: 500}}
      visible={span && open}
      onCancel={handleClose}
      destroyOnClose
      title={
        <div style={{display: 'flex', justifyContent: 'space-between', marginRight: 36}}>
          <Typography.Title level={5}>{assertion ? 'Edit Assertion' : 'Create New Assertion'}</Typography.Title>
          <Typography.Text data-cy="affected-spans-count">
            {affectedSpanMessage(affectedSpanList.length)}
          </Typography.Text>
        </div>
      }
      onOk={form?.submit}
      okButtonProps={{
        type: 'default',
        id: 'add-assertion-modal-ok-button',
      }}
      okText="Save"
    >
      <CreateAssertionForm
        affectedSpanList={affectedSpanList}
        defaultAttributeList={defaultAttributeList}
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
