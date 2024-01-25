import IntegrationIcon from 'components/IntegrationIcon/';
import {GITHUB_ISSUES_URL} from 'constants/Common.constants';
import {IIntegration} from 'constants/Integrations.constants';
import * as S from './TriggerTypeModal.styled';

interface IProps {
  integration: IIntegration;
}

const IntegrationTypeCard = ({integration: {name, title, isActive, isAvailable, url}}: IProps) => (
  <a href={url} target="__blank" aria-disabled={!isActive}>
    <S.IntegrationCardContainer $isActive={isActive} $isSelected={false}>
      <IntegrationIcon integrationName={name} />

      <S.CardContent>
        <div>
          <S.CardTitle $isActive={isActive}>
            {title} {!isAvailable && '*'}{' '}
          </S.CardTitle>
          {!isActive && (
            <S.CardTitle $isActive>
              &nbsp;-{' '}
              <a href={GITHUB_ISSUES_URL} target="_blank">
                Coming soon!
              </a>
            </S.CardTitle>
          )}
        </div>
      </S.CardContent>
    </S.IntegrationCardContainer>
  </a>
);

export default IntegrationTypeCard;
