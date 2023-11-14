import {VariableDefinition} from 'postman-collection';
import {useCallback} from 'react';
import PostmanServiceService, {RequestDefinitionExtended} from 'services/Importers/Postman.service';
import {IPostmanValues, TDraftTestForm} from 'types/Test.types';

export function useSelectTestCallback(
  form: TDraftTestForm<IPostmanValues>,
  requests: RequestDefinitionExtended[],
  variables: VariableDefinition[]
) {
  return useCallback(
    (identifier: string) => {
      PostmanServiceService.updateForm(requests, variables, identifier, form);
    },
    [form, requests, variables]
  );
}
