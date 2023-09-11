import {FormInstance} from 'antd';
import {debounce} from 'lodash';
import {useEffect, useMemo, useState} from 'react';

import {SupportedEditors} from 'constants/Editor.constants';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useAppSelector} from 'redux/hooks';
import SpanSelectors from 'selectors/Span.selectors';
import EditorService from 'services/Editor.service';
import SelectorSuggestionsService from 'services/SelectorSuggestions/SelectorSuggestions.service';
import SpanService from 'services/Span.service';
import useAssertionFormValues from './useAssertionFormValues';
import {IValues} from '../TestSpecForm';

interface IDebounceProps {
  q: string;
  rId: number;
  tId: string;
}

interface IProps {
  form: FormInstance<IValues>;
  runId: number;
  testId: string;
  onValidSelector(isValid: boolean): void;
}

const useQuerySelector = ({form, runId, testId, onValidSelector}: IProps) => {
  const {onClearMatchedSpans, selectedSpan, onTriggerSelector, triggerSelectorResult, isTriggerSelectorError} =
    useSpan();
  const {setSelectorSuggestions} = useTestSpecs();
  const [isLoading, setIsLoading] = useState(true);
  const {currentSelector} = useAssertionFormValues(form);
  const [isValid, setIsValid] = useState(!isTriggerSelectorError);
  const selectedParentSpan = useAppSelector(state =>
    SpanSelectors.selectSpanById(state, selectedSpan?.parentId ?? '', testId, runId)
  );

  const handleSelector = useMemo(
    () =>
      debounce(async ({q, tId, rId}: IDebounceProps) => {
        const isValidSelector = EditorService.getIsQueryValid(SupportedEditors.Selector, q || '');

        setIsValid(isValidSelector);

        if (isValidSelector) await onTriggerSelector(q, tId, rId);

        setIsLoading(false);
      }, 500),
    [onTriggerSelector]
  );

  useEffect(() => {
    setIsLoading(true);
    handleSelector({q: currentSelector, tId: testId, rId: runId});
  }, [handleSelector, currentSelector, runId, testId]);

  useEffect(() => {
    if (!triggerSelectorResult) return;

    const selectedSpanId = selectedSpan?.id ?? '';
    const selectedSpanSelector = SpanService.getSelectorInformation(selectedSpan!);
    const selectedParentSpanSelector = selectedParentSpan ? SpanService.getSelectorInformation(selectedParentSpan) : '';

    const selectorSuggestions = SelectorSuggestionsService.getSuggestions(
      triggerSelectorResult.selector,
      triggerSelectorResult.spanIds,
      selectedSpanId,
      selectedSpanSelector,
      selectedParentSpanSelector
    );
    setSelectorSuggestions(selectorSuggestions);
  }, [selectedParentSpan, selectedSpan, setSelectorSuggestions, triggerSelectorResult]);

  useEffect(() => {
    return () => {
      onClearMatchedSpans();
    };
  }, [onClearMatchedSpans]);

  useEffect(() => {
    setIsValid(!isTriggerSelectorError);
  }, [isTriggerSelectorError]);

  useEffect(() => {
    form.setFields([
      {
        name: 'selector',
        errors: !isValid ? ['Invalid selector'] : [],
      },
    ]);
    onValidSelector(isValid);
  }, [form, isValid, onValidSelector]);

  return {
    spanIdList: triggerSelectorResult?.spanIds ?? [],
    isValid,
    isLoading,
  };
};

export default useQuerySelector;
