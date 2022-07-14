import {FormInstance} from 'antd';
import {Dispatch, SetStateAction, useEffect, useState} from 'react';
import {IUploadCollectionValues} from '../UploadCollection';

export function useValidateFormEffect(
  form: FormInstance<IUploadCollectionValues>
): [boolean, Dispatch<SetStateAction<boolean>>] {
  const [isFormValid, setIsFormValid] = useState(false);
  useEffect(() => {
    async function fn() {
      try {
        await form.validateFields();
        setIsFormValid(true);
      } catch (err) {
        setIsFormValid(false);
      }
    }

    fn();
  }, [form, setIsFormValid]);
  return [isFormValid, setIsFormValid];
}
