import {Operation} from 'components/AllowButton';
import * as S from '../TestSuites/TestSuites.styled';

interface IProps {
  onCreate(): void;
  title?: string;
}

const CreateButton = ({onCreate, title}: IProps) => (
  <S.ActionContainer>
    <S.CreateTestButton operation={Operation.Edit} type="primary" data-cy="create-button" onClick={onCreate}>
      {title || 'Create'}
    </S.CreateTestButton>
  </S.ActionContainer>
);

export default CreateButton;
