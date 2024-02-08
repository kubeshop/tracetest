import {useCallback} from 'react';
import {noop} from 'lodash';
import {Completion, CompletionContext} from '@codemirror/autocomplete';
import {useAppSelector} from 'redux/hooks';
import SpanSelectors from 'selectors/Span.selectors';
import AssertionSelectors from 'selectors/Assertion.selectors';
import VariableSetSelectors from 'selectors/VariableSet.selectors';
import EditorService from 'services/Editor.service';
import {SupportedEditors} from 'constants/Editor.constants';

interface IProps {
  testId: string;
  runId: number;
  onSelect?(option: Completion): void;
  autocompleteCustomValues: string[];
}

const useAutoComplete = ({testId, runId, onSelect = noop, autocompleteCustomValues}: IProps) => {
  const attributeList = useAppSelector(state => {
    const spanIds = SpanSelectors.selectMatchedSpans(state);
    return AssertionSelectors.selectAllUniqueAttributeList(state, testId, runId, spanIds);
  });
  const varEntryList = useAppSelector(state => VariableSetSelectors.selectSelectedVariableSetValues(state, true));

  return useCallback(
    async (context: CompletionContext) =>
      EditorService.getAutocomplete({
        type: SupportedEditors.Expression,
        context,
        attributeList,
        varEntryList,
        customValueList: autocompleteCustomValues,
        onSelect,
      }),
    [attributeList, autocompleteCustomValues, onSelect, varEntryList]
  );
};

export default useAutoComplete;
