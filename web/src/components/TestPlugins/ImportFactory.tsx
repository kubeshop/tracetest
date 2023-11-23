import {ImportTypes} from 'constants/Test.constants';
import {TDraftTestForm} from 'types/Test.types';
import Postman from './Imports/Postman';
import Curl from './Imports/Curl';
import useShortcut from './hooks/useShortcut';
import Definition from './Imports/Definition';

const ImportFactoryMap = {
  [ImportTypes.postman]: Postman,
  [ImportTypes.curl]: Curl,
  [ImportTypes.definition]: Definition,
};

export interface IFormProps {
  form: TDraftTestForm;
}

interface IProps {
  type: ImportTypes;
}

const ImportFactory = ({type}: IProps) => {
  useShortcut();
  const Component = ImportFactoryMap[type];

  return <Component />;
};

export default ImportFactory;
