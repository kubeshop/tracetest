import {IPlugin} from 'types/Plugins.types';
import {GITHUB_ISSUES_URL} from '../../constants/Common.constants';
import * as S from './PluginCard.styled';

interface IProps {
  plugin: IPlugin;
  isSelected: boolean;
  onSelect(plugin: IPlugin): void;
}

const PluginCard = ({plugin: {title, description, isActive}, plugin, onSelect, isSelected}: IProps) => {
  return (
    <S.PluginCard $isSelected={isSelected} $isActive={isActive} onClick={() => isActive && onSelect(plugin)}>
      <S.Circle>{isSelected && <S.Check />}</S.Circle>
      <S.Content>
        <S.Title>
          {title}{' '}
          {!isActive && (
            <>
              - <a href={GITHUB_ISSUES_URL}>Coming soon!</a>
            </>
          )}
        </S.Title>
        <S.Description>{description}</S.Description>
      </S.Content>
    </S.PluginCard>
  );
};

export default PluginCard;
