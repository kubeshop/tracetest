import {useCallback} from 'react';
import {LoadingOutlined} from '@ant-design/icons';
import SpanService from 'services/Span.service';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import SpanSelectors from 'selectors/Span.selectors';
import {useSpan} from 'providers/Span/Span.provider';
import TestOutput from 'models/TestOutput.model';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import {useConfirmationModal} from 'providers/ConfirmationModal/ConfirmationModal.provider';
import {useTestSpecForm} from '../TestSpecForm/TestSpecForm.provider';
import * as S from './CurrentSpanSelector.styled';

interface IProps {
  spanId: string;
}

const CurrentSpanSelector = ({spanId}: IProps) => {
  const {open, isOpen: isTestSpecFormOpen} = useTestSpecForm();
  const {onOpen} = useTestOutput();
  const {onOpen: onConfirmModalOpen} = useConfirmationModal();
  const {
    run: {id: runId},
  } = useTestRun();
  const {
    test: {id: testId},
  } = useTest();
  const {isTriggerSelectorLoading} = useSpan();
  const span = useAppSelector(state => SpanSelectors.selectSpanById(state, spanId, testId, runId));

  const handleOnClick = useCallback(() => {
    const selector = SpanService.getSelectorInformation(span!);
    onConfirmModalOpen({
      title: 'Unsaved changes',
      heading: 'Discard unsaved changes?',
      okText: 'Discard',
      onConfirm: () => {
        if (isTestSpecFormOpen) {
          open({
            isEditing: false,
            selector,
            defaultValues: {
              selector,
            },
          });
        } else {
          onOpen(TestOutput({selector: {query: selector}}));
        }
      },
    });
  }, [isTestSpecFormOpen, onConfirmModalOpen, onOpen, open, span]);

  return (
    <S.Container className="matched">
      <S.FloatingText onClick={() => !isTriggerSelectorLoading && handleOnClick()}>
        {isTriggerSelectorLoading ? (
          <>
            Updating selected span <LoadingOutlined />
          </>
        ) : (
          <>Select as Current span</>
        )}
      </S.FloatingText>
    </S.Container>
  );
};

export default CurrentSpanSelector;
