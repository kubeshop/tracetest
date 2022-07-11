import {UploadFile} from 'antd/es/upload/interface';
import {Collection, VariableDefinition} from 'postman-collection';
import {Dispatch, SetStateAction, useCallback} from 'react';
import {getRequestsFromCollection, RequestDefinitionExtended} from './getRequestsFromCollection';

export function useUploadCollectionCallback(setState: Dispatch<SetStateAction<State>>): (file?: UploadFile) => void {
  return useCallback(
    (file?: UploadFile) => {
      console.log('here');
      const readFile = new FileReader();
      readFile.onload = (e: ProgressEvent<FileReader>) => {
        const contents = e?.target?.result;
        if (contents && typeof contents === 'string') {
          try {
            const collection = new Collection(JSON.parse(contents));
            setState({variables: collection.variables.all(), requests: getRequestsFromCollection(collection)});
          } catch (r) {
            console.error('erro');
          }
        }
      };
      readFile.readAsText(file as any);
    },
    [setState]
  );
}

export interface State {
  requests: RequestDefinitionExtended[];
  variables: VariableDefinition[];
}
