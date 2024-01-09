import {Typography} from 'antd';
import {FramedCodeBlock} from 'components/CodeBlock';
import {CypressCodeSnippet} from 'constants/Automate.constants';
import Test from 'models/Test.model';
import * as S from './Cypress.styled';
import {IMethodChildrenProps} from '../../RunDetailAutomateMethods';

interface IProps extends IMethodChildrenProps {
  test: Test;
}

const Cypress = ({test}: IProps) => (
  <S.Container>
    <S.TitleContainer>
      <S.Title>Cypress Integration</S.Title>
    </S.TitleContainer>
    <Typography.Paragraph>The code snippet below enables you to run this test via a Cypress run.</Typography.Paragraph>
    <FramedCodeBlock title="Cypress code snippet:" language="javascript" value={CypressCodeSnippet(test.name)} />
  </S.Container>
);

export default Cypress;
