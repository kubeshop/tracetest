import {useCallback} from 'react';
import {noop, uniqBy} from 'lodash';
import {Completion, CompletionContext} from '@codemirror/autocomplete';
import {useAppStore} from 'redux/hooks';
import AssertionSelectors from 'selectors/Assertion.selectors';
import EnvironmentSelectors from 'selectors/Environment.selectors';
import SpanSelectors from 'selectors/Span.selectors';
import EditorService from 'services/Editor.service';
import {SupportedEditors} from 'constants/Editor.constants';

interface IProps {
  testId: string;
  runId: string;
  onSelect?(option: Completion): void;
  autocompleteCustomValues: string[];
}

const useAutoComplete = ({testId, runId, onSelect = noop, autocompleteCustomValues}: IProps) => {
  const {getState} = useAppStore();

  const getAttributeList = useCallback(() => {
    const state = getState();
    const spanIdList = SpanSelectors.selectMatchedSpans(state);
    const attributeList = AssertionSelectors.selectAttributeList(state, testId, runId, spanIdList);

    return uniqBy(attributeList, 'key');
  }, [getState, runId, testId]);

  const getSelectedEnvironmentEntryList = useCallback(() => {
    const state = getState();

    return EnvironmentSelectors.selectSelectedEnvironmentValues(state, true);
  }, [getState]);

  return useCallback(
    async (context: CompletionContext) => {
      const attributeList = getAttributeList();
      const envEntryList = getSelectedEnvironmentEntryList();

      return EditorService.getAutocomplete({
        type: SupportedEditors.Expression,
        context,
        attributeList,
        envEntryList,
        customValueList: autocompleteCustomValues,
        onSelect,
      });
    },
    [autocompleteCustomValues, getAttributeList, getSelectedEnvironmentEntryList, onSelect]
  );
};

export default useAutoComplete;
