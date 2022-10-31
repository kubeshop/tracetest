import {useCallback} from 'react';
import {CompletionContext} from '@codemirror/autocomplete';
import EditorService from 'services/Editor.service';
import {SupportedEditors} from 'constants/Editor.constants';
import {useAppStore} from 'redux/hooks';
import EnvironmentSelectors from 'selectors/Environment.selectors';

const useAutoComplete = () => {
  const {getState} = useAppStore();

  const getSelectedEnvironmentEntryList = useCallback(() => {
    const state = getState();

    return EnvironmentSelectors.selectSelectedEnvironmentValues(state);
  }, [getState]);

  return useCallback(
    async (context: CompletionContext) => {
      const envEntryList = getSelectedEnvironmentEntryList();

      return EditorService.getAutocomplete({type: SupportedEditors.Interpolation, context, envEntryList});
    },
    [getSelectedEnvironmentEntryList]
  );
};

export default useAutoComplete;
