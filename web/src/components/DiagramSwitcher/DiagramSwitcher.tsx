import {SupportedDiagrams} from '../Diagram/Diagram';
import React from 'react';
import SearchInput from '../SearchInput';
import * as S from './DiagramSwitcher.styled';

interface IProps {
  onSearch(search: string): void;

  onTypeChange(type: SupportedDiagrams): void;

  selectedType: SupportedDiagrams;
}

const DiagramSwitcher: React.FC<IProps> = ({onSearch, onTypeChange, selectedType}) => {
  return (
    <S.DiagramSwitcher>
      <S.Switch>
        <S.DAGIcon
          $isSelected={selectedType === SupportedDiagrams.DAG}
          onClick={() => onTypeChange(SupportedDiagrams.DAG)}
        />
        <S.TimelineIcon
          $isSelected={selectedType === SupportedDiagrams.Timeline}
          onClick={() => onTypeChange(SupportedDiagrams.Timeline)}
        />
      </S.Switch>
      <SearchInput onSearch={onSearch} width="100%" placeholder="Search in trace (Not implemented yet)" />
    </S.DiagramSwitcher>
  );
};

export default DiagramSwitcher;
