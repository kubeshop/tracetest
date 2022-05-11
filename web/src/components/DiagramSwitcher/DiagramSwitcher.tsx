import {debounce} from 'lodash';
import {useMemo} from 'react';
import {SupportedDiagrams} from '../Diagram/Diagram';
import * as S from './DiagramStories.styled';

interface IDiagramSwitcherProps {
  onSearch(search: string): void;
  onTypeChange(type: SupportedDiagrams): void;
  selectedType: SupportedDiagrams;
}

const DiagramSwitcher: React.FC<IDiagramSwitcherProps> = ({onSearch, onTypeChange, selectedType}) => {
  const handleSearch = useMemo(
    () =>
      debounce(event => {
        onSearch(event.target.value);
      }, 500),
    [onSearch]
  );

  return (
    <S.DiagramSwitcher>
      <S.Switch>
        <S.DAGIcon
          isSelected={selectedType === SupportedDiagrams.DAG}
          onClick={() => onTypeChange(SupportedDiagrams.DAG)}
        />
        <S.TimelineIcon
          isSelected={selectedType === SupportedDiagrams.Timeline}
          onClick={() => onTypeChange(SupportedDiagrams.Timeline)}
        />
      </S.Switch>
      <S.SearchInput onChange={handleSearch} />
    </S.DiagramSwitcher>
  );
};

export default DiagramSwitcher;
