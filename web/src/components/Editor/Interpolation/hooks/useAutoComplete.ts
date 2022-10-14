import {useCallback} from 'react';
import {CompletionContext} from '@codemirror/autocomplete';
import EditorService from 'services/Editor.service';
import {SupportedEditors} from 'constants/Editor.constants';

const useAutoComplete = () => {
  return useCallback(async (context: CompletionContext) => {
    return EditorService.getAutocomplete(SupportedEditors.Interpolation, context);
  }, []);
};

export default useAutoComplete;
