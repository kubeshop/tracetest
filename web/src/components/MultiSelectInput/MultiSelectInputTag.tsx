import {CustomTagProps} from 'antd/node_modules/rc-select/lib/BaseSelect';
import {SEPARATOR} from './MultiSelectInput';
import * as S from './MultiSelectInput.styled';

interface IProps extends CustomTagProps {
  onDeselect(entryNumber: number): void;
  entryListCount: number;
}

const MultiSelectInputTag: React.FC<IProps> = ({value, entryListCount, onDeselect, ...props}) => {
  const [, stepNumber, entryNumber] = value.split(SEPARATOR);
  const isLast = Number(stepNumber) === entryListCount - 1;

  return (
    <S.SelectedTag
      {...props}
      closable={isLast}
      onMouseDown={e => {
        if (isLast) e.stopPropagation();
      }}
      onClose={() => {
        onDeselect(Number(entryNumber));
      }}
      isLast={isLast}
    >
      {props.label}
    </S.SelectedTag>
  );
};

export default MultiSelectInputTag;
