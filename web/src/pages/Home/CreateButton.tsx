import {Operation} from 'components/AllowButton';
import * as S from '../TestSuites/TestSuites.styled';

interface IProps {
  onCreate(): void;
  title?: string;
  dataCy?: string;
}

const CreateButton = ({onCreate, title, dataCy}: IProps) => (
  <S.ActionContainer>
    <S.CreateTestButton operation={Operation.Edit} type="primary" data-cy={dataCy} onClick={onCreate}>
      {title || 'Create'}
    </S.CreateTestButton>
  </S.ActionContainer>
);

export default CreateButton;
