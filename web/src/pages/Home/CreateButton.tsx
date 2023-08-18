import * as S from '../TestSuites/TestSuites.styled';

interface IProps {
  onCreate(): void;
}

const CreateButton = ({onCreate}: IProps) => {
  return (
    <S.ActionContainer>
      <S.CreateTestButton type="primary" data-cy="create-button" onClick={onCreate}>
        Create
      </S.CreateTestButton>
    </S.ActionContainer>
  );
};

export default CreateButton;
