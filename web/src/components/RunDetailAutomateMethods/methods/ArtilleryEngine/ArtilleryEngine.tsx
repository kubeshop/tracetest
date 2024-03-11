import {ArtilleryEngineCodeSnippet} from 'constants/Automate.constants';
import {ReadOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import {FramedCodeBlock} from 'components/CodeBlock';
import Test from 'models/Test.model';
import * as S from './ArtilleryEngine.styled';
import {IMethodChildrenProps} from '../../RunDetailAutomateMethods';

interface IProps extends IMethodChildrenProps {
  test: Test;
}

const ARTILLERY_DOCS = 'https://docs.tracetest.io/tools-and-integrations/artillery-engine';

const ArtilleryEngine = ({test}: IProps) => (
  <S.Container>
    <S.TitleContainer>
      <S.Title>Artillery Engine Integration</S.Title>
      <a href={ARTILLERY_DOCS} target="_blank">
        <ReadOutlined />
      </a>
    </S.TitleContainer>
    <Typography.Paragraph>
      The code snippet below enables you to run this test via a Artillery run.
    </Typography.Paragraph>
    <FramedCodeBlock
      title="Install:"
      language="bash"
      value="npm i -g artillery-engine-tracetest"
      minHeight="50px"
      maxHeight="50px"
    />
    <br />
    <FramedCodeBlock
      title="Artillery Test Script:"
      language="yaml"
      value={ArtilleryEngineCodeSnippet(test.id)}
    />
  </S.Container>
);

export default ArtilleryEngine;
