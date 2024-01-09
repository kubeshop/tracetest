import {Typography} from 'antd';
import {FramedCodeBlock} from 'components/CodeBlock';
import {PlaywrightCodeSnippet} from 'constants/Automate.constants';
import Test from 'models/Test.model';
import * as S from './Playwright.styled';
import {IMethodChildrenProps} from '../../RunDetailAutomateMethods';

interface IProps extends IMethodChildrenProps {
  test: Test;
}

const Playwright = ({test}: IProps) => (
  <S.Container>
    <S.TitleContainer>
      <S.Title>Playwright Integration</S.Title>
    </S.TitleContainer>
    <Typography.Paragraph>The code snippet below enables you to run this test via a Playwright run.</Typography.Paragraph>
    <FramedCodeBlock title="Cypress code snippet:" language="javascript" value={PlaywrightCodeSnippet(test.name)} />
  </S.Container>
);

export default Playwright;
