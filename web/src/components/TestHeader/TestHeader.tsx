import ResourceCardActions from 'components/ResourceCard/ResourceCardActions';
import {Link} from 'react-router-dom';
import * as S from './TestHeader.styled';

interface IProps {
  description: string;
  id: string;
  canEdit: boolean;
  onEdit(): void;
  onDelete(): void;
  title: string;
  runButton: React.ReactElement;
}

const TestHeader = ({description, id, canEdit, onEdit, onDelete, title, runButton}: IProps) => (
  <S.Container $isWhite>
    <S.Section>
      <Link to="/" data-cy="test-header-back-button">
        <S.BackIcon />
      </Link>
      <div>
        <S.Title data-cy="test-details-name">{title}</S.Title>
        <S.Text>{description}</S.Text>
      </div>
    </S.Section>
    <S.Section>
      {runButton}
      <ResourceCardActions id={id} onDelete={onDelete} onEdit={onEdit} canEdit={canEdit} />
    </S.Section>
  </S.Container>
);

export default TestHeader;
