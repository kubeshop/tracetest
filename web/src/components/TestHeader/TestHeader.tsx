import ResourceCardActions from 'components/ResourceCard/ResourceCardActions';
import * as S from './TestHeader.styled';

interface IProps {
  description: string;
  id: string;
  onBack(): void;
  onDelete(): void;
  title: string;
}

const TestHeader = ({description, id, onBack, onDelete, title}: IProps) => (
  <S.Container $isWhite>
    <S.Section>
      <S.BackIcon data-cy="test-header-back-button" onClick={onBack} />
      <div>
        <S.Title data-cy="test-details-name">{title}</S.Title>
        <S.Text>{description}</S.Text>
      </div>
    </S.Section>
    <S.Section>
      <ResourceCardActions id={id} onDelete={onDelete} />
    </S.Section>
  </S.Container>
);

export default TestHeader;
