import TestCardActions from 'components/TestCard/TestCardActions';
import {useMenuDeleteCallback} from 'pages/Home/useMenuDeleteCallback';
import {TTest} from 'types/Test.types';
import * as S from './TestHeader.styled';

interface IProps {
  onBack(): void;
  test: TTest;
}

const TestHeader = ({onBack, test: {id, name, trigger, version = 1}, test}: IProps) => {
  const onDelete = useMenuDeleteCallback();

  return (
    <S.Container>
      <S.Section>
        <S.BackIcon data-cy="test-header-back-button" onClick={onBack} />
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
