import ResourceCardActions from 'components/ResourceCard/ResourceCardActions';
import * as S from './TestHeader.styled';
import VariableSetSelector from '../VariableSetSelector/VariableSetSelector';

interface IProps {
  description: string;
  id: string;
  shouldEdit: boolean;
  onBack(): void;
  onEdit(): void;
  onDelete(): void;
  onDuplicate(): void;
  title: string;
  runButton: React.ReactNode;
}

const TestHeader = ({description, id, shouldEdit, onBack, onEdit, onDelete, onDuplicate, title, runButton}: IProps) => (
  <S.Container $isWhite>
    <S.Section>
      <a onClick={onBack} data-cy="test-header-back-button">
        <S.BackIcon />
      </a>
      <div>
        <S.Title data-cy="test-details-name">{title}</S.Title>
        <S.Text>{description}</S.Text>
      </div>
    </S.Section>
    <S.Section>
      <VariableSetSelector />
      {runButton}
      <ResourceCardActions
        id={id}
        onDuplicate={onDuplicate}
        onDelete={onDelete}
        onEdit={onEdit}
        shouldEdit={shouldEdit}
      />
    </S.Section>
  </S.Container>
);

export default TestHeader;
