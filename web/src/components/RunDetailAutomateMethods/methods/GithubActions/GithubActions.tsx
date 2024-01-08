import {Typography} from 'antd';
import {getServerBaseUrl} from 'utils/Common';
import {GITHUB_ACTION_URL, CLI_DOCS_URL} from 'constants/Common.constants';
import {withCustomization} from 'providers/Customization';
import {FramedCodeBlock} from 'components/CodeBlock';
import * as S from './GithubActions.styled';

const actionConfig = `- name: Configure Tracetest CLI
  uses: kubeshop/tracetest-github-action@v1
  with:
    endpoint: ${getServerBaseUrl()}
`;

const GithubActions = () => {
  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>GitHub Action Configuration</S.Title>
      </S.TitleContainer>
      <Typography.Paragraph>
        Integrate Tracetest into your GitHub pipeline by adding this snippet to your workflow steps to utilize{' '}
        <a href={CLI_DOCS_URL} target="__blank">
          Tracetest CLI
        </a>{' '}
        for test runs:
      </Typography.Paragraph>
      <FramedCodeBlock title="Snippet" language="yaml" value={actionConfig} minHeight="120px" maxHeight="120px" />
      <S.Subtitle type="secondary">
        The endpoint parameter is the base address where your Tracetest Core Server is installed. <br /> Here&apos;s a
        full example of how to use it:{' '}
        <a href={GITHUB_ACTION_URL} target="__blank">
          tracetest-cli-with-tracetest-core.yml
        </a>
      </S.Subtitle>
    </S.Container>
  );
};

export default withCustomization(GithubActions, 'githubActions');
