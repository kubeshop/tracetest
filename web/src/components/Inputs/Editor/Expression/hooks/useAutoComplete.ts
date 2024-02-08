import {useCallback} from 'react';
import {noop} from 'lodash';
import {Completion, CompletionContext} from '@codemirror/autocomplete';
import {useAppStore} from 'redux/hooks';
import VariableSetSelectors from 'selectors/VariableSet.selectors';
import {selectExpressionAttributeList} from 'selectors/Editor.selectors';
import EditorService from 'services/Editor.service';
import {SupportedEditors} from 'constants/Editor.constants';

interface IProps {
  testId: string;
  runId: number;
  onSelect?(option: Completion): void;
  autocompleteCustomValues: string[];
}

const useAutoComplete = ({testId, runId, onSelect = noop, autocompleteCustomValues}: IProps) => {
  const {getState} = useAppStore();

  const getAttributeList = useCallback(
    () => selectExpressionAttributeList(getState(), testId, runId),
    [getState, runId, testId]
  );

  const getSelectedVariableSetEntryList = useCallback(() => {
    const state = getState();

    return VariableSetSelectors.selectSelectedVariableSetValues(state, true);
  }, [getState]);

  return useCallback(
    async (context: CompletionContext) => {
      const attributeList = getAttributeList();
      const varEntryList = getSelectedVariableSetEntryList();

      return EditorService.getAutocomplete({
        type: SupportedEditors.Expression,
        context,
        attributeList,
        varEntryList,
        customValueList: autocompleteCustomValues,
        onSelect,
      });
    },
    [autocompleteCustomValues, getAttributeList, getSelectedVariableSetEntryList, onSelect]
  );
};

export default useAutoComplete;
