import {Typography} from 'antd';
import {ReadOutlined} from '@ant-design/icons';
import {FramedCodeBlock} from 'components/CodeBlock';
import {CypressCodeSnippet} from 'constants/Automate.constants';
import Test from 'models/Test.model';
import * as S from './Cypress.styled';
import {IMethodChildrenProps} from '../../RunDetailAutomateMethods';

interface IProps extends IMethodChildrenProps {
  test: Test;
}

const CYPRESS_DOCS = 'https://docs.tracetest.io/tools-and-integrations/cypress';

const Cypress = ({test}: IProps) => (
  <S.Container>
    <S.TitleContainer>
      <S.Title>Cypress Integration</S.Title>
      <a href={CYPRESS_DOCS} target="_blank">
        <ReadOutlined />
      </a>
    </S.TitleContainer>
    <Typography.Paragraph>The code snippet below enables you to run this test via a Cypress run.</Typography.Paragraph>
    <FramedCodeBlock
      title="Install:"
      language="bash"
      value="npm i -g @tracetest/cypress"
      minHeight="50px"
      maxHeight="50px"
    />
    <br />
    <FramedCodeBlock title="Cypress code snippet:" language="javascript" value={CypressCodeSnippet(test.name)} />
  </S.Container>
);

export default Cypress;
