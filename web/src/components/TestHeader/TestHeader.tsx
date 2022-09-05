import {useNavigate} from 'react-router-dom';
import TestCardActions from 'components/TestCard/TestCardActions';
import {TTest} from 'types/Test.types';
import useDeleteTest from 'hooks/useDeleteTest';
import * as S from './TestHeader.styled';

interface IProps {
  test: TTest;
}

const TestHeader = ({test: {id, name, trigger, version = 1}, test}: IProps) => {
  const onDelete = useDeleteTest();
  const navigate = useNavigate();

  return (
    <S.Container>
      <S.Section>
        <S.BackIcon data-cy="test-header-back-button" onClick={() => navigate(-1)} />
        <div>
          <S.Title data-cy="test-details-name">
            {name} (v{version})
          </S.Title>
          <S.Text>
            {trigger.type.toUpperCase()} - {trigger.method.toUpperCase()} - {trigger.entryPoint}
          </S.Text>
        </div>
      </S.Section>
      <S.Section>
        <TestCardActions testId={id} onDelete={() => onDelete(test)} />
      </S.Section>
    </S.Container>
  );
};

export default TestHeader;
