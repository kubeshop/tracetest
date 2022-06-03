import {Tooltip} from 'antd';
import React from 'react';
import GuidedTourService, {GuidedTours} from '../../services/GuidedTour.service';
import {SupportedDiagrams} from '../Diagram/Diagram';
import {Steps} from '../GuidedTour/traceStepList';
import SearchInput from '../SearchInput';
import * as S from './DiagramSwitcher.styled';

interface IProps {
  onSearch(search: string): void;
  onTypeChange(type: SupportedDiagrams): void;
  selectedType: SupportedDiagrams;
}

const color = '#FFFFFF';

const DiagramSwitcher: React.FC<IProps> = ({onSearch, onTypeChange, selectedType}) => {
  return (
    <S.DiagramSwitcher>
      <S.Switch data-tour={GuidedTourService.getStep(GuidedTours.Trace, Steps.Switcher)}>
        <Tooltip title="Graph view" color={color}>
          <S.DAGIcon
            $isSelected={selectedType === SupportedDiagrams.DAG}
            onClick={() => onTypeChange(SupportedDiagrams.DAG)}
          />
        </Tooltip>
        <Tooltip title="Timeline view" color={color}>
          <S.TimelineIcon
            $isSelected={selectedType === SupportedDiagrams.Timeline}
            onClick={() => onTypeChange(SupportedDiagrams.Timeline)}
          />
        </Tooltip>
      </S.Switch>
      <SearchInput onSearch={onSearch} width="100%" placeholder="Search in trace (Not implemented yet)" />
    </S.DiagramSwitcher>
  );
};

export default DiagramSwitcher;
