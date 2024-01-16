import {GITHUB_ISSUES_URL} from 'constants/Common.constants';
import {IPlugin} from 'types/Plugins.types';
import * as S from './TriggerTypeModal.styled';

interface IProps {
  onClick(plugin: IPlugin): void;
  plugin: IPlugin;
  isSelected: boolean;
}

const TriggerTypeCard = ({plugin: {name, title, description, isActive}, plugin, onClick, isSelected}: IProps) => (
  <S.CardContainer
    data-cy={`${name.toLowerCase()}-plugin`}
    onClick={() => isActive && onClick(plugin)}
    $isActive={isActive}
    $isSelected={isSelected}
  >
    <S.Circle $isActive={isActive}>
      <S.Check className="check" />
    </S.Circle>

    <S.CardContent>
      <div>
        <S.CardTitle $isActive={isActive}>{title} </S.CardTitle>
        {!isActive && (
          <S.CardTitle $isActive>
            &nbsp;-{' '}
            <a href={GITHUB_ISSUES_URL} target="_blank">
              Coming soon!
            </a>
          </S.CardTitle>
        )}
      </div>
      <S.CardDescription $isActive={isActive}>{description}</S.CardDescription>
    </S.CardContent>
  </S.CardContainer>
);

export default TriggerTypeCard;
