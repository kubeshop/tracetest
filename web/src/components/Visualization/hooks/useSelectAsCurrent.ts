import {useCallback} from 'react';
import {useTestSpecForm} from 'components/TestSpecForm/TestSpecForm.provider';
import Span from 'models/Span.model';
import TestOutput from 'models/TestOutput.model';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestOutput} from 'providers/TestOutput/TestOutput.provider';
import SpanService from 'services/Span.service';

interface IProps {
  selected: boolean;
  matched: boolean;
  span: Span;
}

interface IUseSelectAsCurrent {
  isLoading: boolean;
  onSelectAsCurrent(): void;
  showSelectAsCurrent: boolean;
}

const useSelectAsCurrent = ({selected, matched, span}: IProps): IUseSelectAsCurrent => {
  const {isTriggerSelectorLoading, matchedSpans} = useSpan();
  const {isOpen: isTestSpecFormOpen, open: onOpenTestSpecForm} = useTestSpecForm();
  const {isOpen: isTestOutputFormOpen, onOpen: onOpenTestOutputForm} = useTestOutput();

  const showSelectAsCurrent =
    selected && !matched && !!matchedSpans.length && (isTestSpecFormOpen || isTestOutputFormOpen);

  const onSelectAsCurrent = useCallback(() => {
    const selector = SpanService.getSelectorInformation(span);
    if (isTestSpecFormOpen) {
      onOpenTestSpecForm({
        isEditing: false,
        selector,
        defaultValues: {selector},
      });
    } else {
      onOpenTestOutputForm(TestOutput({selector, selectorParsed: {query: selector}}));
    }
  }, [isTestSpecFormOpen, onOpenTestSpecForm, onOpenTestOutputForm, span]);

  return {isLoading: isTriggerSelectorLoading, onSelectAsCurrent, showSelectAsCurrent};
};

export default useSelectAsCurrent;
