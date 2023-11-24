import {noop} from 'lodash';
import {ImportTypes} from 'constants/Test.constants';
import ImportCard from './ImportCard';
import * as S from './ImportSelector.styled';

const importList = [ImportTypes.definition, ImportTypes.curl, ImportTypes.postman];

interface IProps {
  value?: ImportTypes;
  onChange?(value: ImportTypes): void;
}

const ImportSelector = ({value, onChange = noop}: IProps) => {
  return (
    <S.CardList>
      {importList.map(importType => (
        <ImportCard
          key={importType}
          onClick={() => onChange(importType)}
          isSelected={value === importType}
          name={importType}
        />
      ))}
    </S.CardList>
  );
};

export default ImportSelector;
