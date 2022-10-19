import {FormInstance} from 'antd';
import {debounce} from 'lodash';
import {useEffect, useMemo, useState} from 'react';

import useEditorValidate from 'components/Editor/hooks/useEditorValidate';
import {SupportedEditors} from 'constants/Editor.constants';
import {useSpan} from 'providers/Span/Span.provider';
import {useTestSpecs} from 'providers/TestSpecs/TestSpecs.provider';
import {useLazyGetSelectedSpansQuery} from 'redux/apis/TraceTest.api';
import SelectorSuggestionsService from 'services/SelectorSuggestions/SelectorSuggestions.service';
import useAssertionFormValues from './useAssertionFormValues';
import {IValues} from '../TestSpecForm';

interface IDebounceProps {
  q: string;
  rId: string;
  tId: string;
}

interface IProps {
  form: FormInstance<IValues>;
  runId: string;
  testId: string;
  onValidSelector(isValid: boolean): void;
}

const useQuerySelector = ({form, runId, testId, onValidSelector}: IProps) => {
  const {onSetMatchedSpans, onClearMatchedSpans} = useSpan();
  const {setSelectorSuggestions} = useTestSpecs();
  const {currentSelector} = useAssertionFormValues(form);
  const [onTriggerSelectedSpans, {data: selectedSpans, isError}] = useLazyGetSelectedSpansQuery();
  const [isValid, setIsValid] = useState(!isError);
  const getIsValidSelector = useEditorValidate();

  const handleSelector = useMemo(
    () =>
      debounce(async ({q, tId, rId}: IDebounceProps) => {
        const isValidSelector = getIsValidSelector(SupportedEditors.Selector, q);

        setIsValid(isValidSelector);
        if (isValidSelector) {
          const selectedSpansData = await onTriggerSelectedSpans({
            query: q,
            testId: tId,
            runId: rId,
          }).unwrap();

          onSetMatchedSpans(selectedSpansData.spanIds);
          setSelectorSuggestions(SelectorSuggestionsService.getSuggestions(selectedSpansData.selector));
        }
      }, 500),
    [getIsValidSelector, onSetMatchedSpans, onTriggerSelectedSpans, setSelectorSuggestions]
  );

  useEffect(() => {
    handleSelector({q: currentSelector, tId: testId, rId: runId});
  }, [handleSelector, currentSelector, runId, testId]);

  useEffect(() => {
    return () => {
      onClearMatchedSpans();
    };
  }, [onClearMatchedSpans]);

  useEffect(() => {
    setIsValid(!isError);
  }, [isError]);

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
    spanIdList: selectedSpans?.spanIds ?? [],
    isValid,
  };
};

export default useQuerySelector;
