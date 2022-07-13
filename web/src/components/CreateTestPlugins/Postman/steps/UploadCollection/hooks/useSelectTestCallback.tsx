import {FormInstance} from 'antd';
import {VariableDefinition} from 'postman-collection';
import {useCallback} from 'react';
import PostmanServiceService, {RequestDefinitionExtended} from 'services/PostmanService.service';
import {IRequestDetailsValues} from '../UploadCollection';

export function useSelectTestCallback(
  form: FormInstance<IRequestDetailsValues>,
  setTransientUrl: React.Dispatch<React.SetStateAction<string>>,
  requests: RequestDefinitionExtended[],
  variables: VariableDefinition[]
) {
  return useCallback(
    (identifier: string) => {
      PostmanServiceService.updateForm(requests, variables, identifier, form, setTransientUrl);
    },
    [form, setTransientUrl, requests, variables]
  );
}
