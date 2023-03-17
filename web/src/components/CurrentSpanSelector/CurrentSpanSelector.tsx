import {useCallback} from 'react';
import SpanService from 'services/Span.service';
import {useTest} from 'providers/Test/Test.provider';
import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useAppSelector} from 'redux/hooks';
import SpanSelectors from 'selectors/Span.selectors';
import {useTestSpecForm} from '../TestSpecForm/TestSpecForm.provider';
import * as S from './CurrentSpanSelector.styled';

interface IProps {
  spanId: string;
}

const CurrentSpanSelector = ({spanId}: IProps) => {
  const {open} = useTestSpecForm();
  const {
    run: {id: runId},
  } = useTestRun();
  const {
    test: {id: testId},
  } = useTest();
  const span = useAppSelector(state => SpanSelectors.selectSpanById(state, spanId, testId, runId));

  const handleOnClick = useCallback(() => {
    const selector = SpanService.getSelectorInformation(span!);
    open({
      isEditing: false,
      selector,
      defaultValues: {
        selector,
      },
    });
  }, [open, span]);

  return (
    <S.Container className="matched">
      <S.FloatingText onClick={() => handleOnClick()}>Select as Current span</S.FloatingText>
    </S.Container>
  );
};

export default CurrentSpanSelector;
