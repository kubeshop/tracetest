import {useCallback} from 'react';
import {CompletionContext} from '@codemirror/autocomplete';
import EditorService from 'services/Editor.service';
import {SupportedEditors} from 'constants/Editor.constants';
import {useAppStore} from 'redux/hooks';
import VariableSetSelectors from 'selectors/VariableSet.selectors';

const useAutoComplete = () => {
  const {getState} = useAppStore();

  const getSelectedVariableSetEntryList = useCallback(() => {
    const state = getState();

    return VariableSetSelectors.selectSelectedVariableSetValues(state, false);
  }, [getState]);

  return useCallback(
    async (context: CompletionContext) => {
      const envEntryList = getSelectedVariableSetEntryList();

      return EditorService.getAutocomplete({type: SupportedEditors.Interpolation, context, envEntryList});
    },
    [getSelectedVariableSetEntryList]
  );
};

export default useAutoComplete;
