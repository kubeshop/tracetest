import {Tooltip} from 'antd';
import Environment from 'models/Environment.model';
import {useCallback} from 'react';
import * as S from './EnvironmentSelector.styled';

interface IProps {
  environment: Environment;
  isHovering: boolean;
  onEditClick(environment: Environment): void;
}

const EnvironmentSelectorEntry = ({environment: {name}, environment, isHovering, onEditClick}: IProps) => {
  const handleClick = useCallback(
    (event: React.MouseEvent<HTMLSpanElement>) => {
      event.stopPropagation();
      onEditClick(environment);
    },
    [environment, onEditClick]
  );

  return (
    <S.EnvironmentSelectorEntryContainer>
      {name}
      {isHovering && (
        <Tooltip title="Edit Environment">
          <S.EnvironmentSelectorEntryIcon onClick={handleClick} />
        </Tooltip>
      )}
    </S.EnvironmentSelectorEntryContainer>
  );
};

export default EnvironmentSelectorEntry;
