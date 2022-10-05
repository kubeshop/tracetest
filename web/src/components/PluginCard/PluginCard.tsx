import {IPlugin} from 'types/Plugins.types';
import {GITHUB_ISSUES_URL} from 'constants/Common.constants';
import * as S from './PluginCard.styled';

interface IProps {
  plugin: IPlugin;
  isSelected: boolean;
  onSelect(plugin: IPlugin): void;
}

const PluginCard = ({plugin: {name, title, description, isActive}, plugin, onSelect, isSelected}: IProps) => {
  return (
    <S.PluginCard
      data-cy={`${name.toLowerCase()}-plugin`}
      $isSelected={isSelected}
      $isActive={isActive}
      onClick={() => isActive && onSelect(plugin)}
    >
      <S.Circle $isActive={isActive}>{isSelected && <S.Check />}</S.Circle>
      <S.Content>
        <div>
          <S.Title $isActive={isActive}>{title} </S.Title>
          {!isActive && (
            <S.Title $isActive>
              &nbsp;-{' '}
              <a href={GITHUB_ISSUES_URL} target="_blank">
                Coming soon!
              </a>
            </S.Title>
          )}
        </div>
        <S.Description $isActive={isActive}>{description}</S.Description>
      </S.Content>
    </S.PluginCard>
  );
};

export default PluginCard;
