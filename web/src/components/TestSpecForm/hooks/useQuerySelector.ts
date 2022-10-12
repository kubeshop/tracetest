import {FormInstance} from 'antd';
import {debounce} from 'lodash';
import {useEffect, useMemo, useState} from 'react';
import {SupportedEditors} from 'constants/Editor.constants';
import {useSpan} from 'providers/Span/Span.provider';
import {useLazyGetSelectedSpansQuery} from 'redux/apis/TraceTest.api';
import useEditorValidate from 'components/Editor/hooks/useEditorValidate';
import {IValues} from '../TestSpecForm';
import useAssertionFormValues from './useAssertionFormValues';

interface IDebouceProps {
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
  const {currentSelector} = useAssertionFormValues(form);
  const [onTriggerSelectedSpans, {data: spanIdList = [], isError}] = useLazyGetSelectedSpansQuery();
  const [isValid, setIsValid] = useState(!isError);
  const getIsValidSelector = useEditorValidate();

  const handleSelector = useMemo(
    () =>
      debounce(async ({q, tId, rId}: IDebouceProps) => {
        const isValidSelector = getIsValidSelector(SupportedEditors.Selector, q);

        setIsValid(isValidSelector);
        if (isValidSelector) {
          const idList = await onTriggerSelectedSpans({
            query: q,
            testId: tId,
            runId: rId,
          }).unwrap();

          onSetMatchedSpans(idList);
        }
      }, 500),
    [getIsValidSelector, onSetMatchedSpans, onTriggerSelectedSpans]
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
    spanIdList,
    isValid,
  };
};

export default useQuerySelector;
