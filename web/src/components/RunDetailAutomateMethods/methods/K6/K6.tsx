import {K6CodeSnippet} from 'constants/Automate.constants';
import {ReadOutlined} from '@ant-design/icons';
import {Typography} from 'antd';
import {FramedCodeBlock} from 'components/CodeBlock';
import * as S from './K6.styled';

const K6_DOCS = 'https://docs.tracetest.io/tools-and-integrations/k6';

const K6 = () => (
  <S.Container>
    <S.TitleContainer>
      <S.Title>K6 Integration</S.Title>
      <a href={K6_DOCS} target="_blank">
        <ReadOutlined />
      </a>
    </S.TitleContainer>
    <Typography.Paragraph>The code snippet below enables you to run this test via a K6 run.</Typography.Paragraph>
    <FramedCodeBlock
      title="Bundle:"
      language="bash"
      value="xk6 build v0.49.0 --with github.com/kubeshop/xk6-tracetest"
      minHeight="50px"
      maxHeight="50px"
    />
    <br />
    <FramedCodeBlock
      title="K6 Test Script:"
      minHeight="50px"
      maxHeight="50px"
      language="bash"
      value={K6CodeSnippet()}
    />
  </S.Container>
);

export default K6;
