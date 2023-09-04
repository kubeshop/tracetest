import {Tooltip} from 'antd';
import VariableSet from 'models/VariableSet.model';
import {useCallback} from 'react';
import * as S from './VariableSetSelector.styled';

interface IProps {
  variableSet: VariableSet;
  isAllowed: boolean;
  isHovering: boolean;
  onEditClick(variableSet: VariableSet): void;
}

const VariableSetSelectorEntry = ({variableSet: {name}, variableSet, isAllowed, isHovering, onEditClick}: IProps) => {
  const handleClick = useCallback(
    event => {
      event.stopPropagation();
      onEditClick(variableSet);
    },
    [variableSet, onEditClick]
  );

  return (
    <S.VarsSelectorEntryContainer>
      {name}
      {isAllowed && isHovering && (
        <Tooltip title="Edit Variable Set">
          <S.VarsSelectorEntryIcon onClick={handleClick} />
        </Tooltip>
      )}
    </S.VarsSelectorEntryContainer>
  );
};

export default VariableSetSelectorEntry;
