import * as S from './CreateTestHeader.styled';

interface IProps {
  onBack(): void;
}

const CreateTestHeader = ({onBack}: IProps) => {
  return (
    <S.CreateTestHeader data-cy="create-test-header">
      <S.Content>
        <S.BackIcon data-cy="test-header-back-button" onClick={onBack} />
        <div>
          <S.Row>
            <S.Name data-cy="test-details-name">Create Test</S.Name>
          </S.Row>
        </div>
      </S.Content>
    </S.CreateTestHeader>
  );
};

export default CreateTestHeader;
