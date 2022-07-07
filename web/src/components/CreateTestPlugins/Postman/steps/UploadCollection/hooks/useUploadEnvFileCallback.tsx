import {UploadFile} from 'antd/es/upload/interface';
import {VariableScope} from 'postman-collection';
import {useCallback} from 'react';

export function useUploadEnvFileCallback(): (file?: UploadFile) => void {
  return useCallback((file?: UploadFile) => {
    const readFile = new FileReader();
    readFile.onload = (e: ProgressEvent<FileReader>) => {
      const contents = e?.target?.result;
      if (contents && typeof contents === 'string') {
        try {
          const variableList = new VariableScope(JSON.parse(contents));
          console.log(variableList.values?.values());
        } catch (r: any) {
          alert('erro');
        }
      }
    };
    readFile.readAsText(file as any);
  }, []);
}
