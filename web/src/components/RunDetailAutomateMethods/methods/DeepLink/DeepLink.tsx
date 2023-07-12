import {Typography} from 'antd';
import {FramedCodeBlock} from 'components/CodeBlock';
import * as S from './DeepLink.styled';
import Controls from './Controls';
import {IMethodProps} from '../../RunDetailAutomateMethods';
import useDeepLink from './hooks/useDeepLink';

const DeepLink = ({test, environmentId, run: {environment}}: IMethodProps) => {
  const {deepLink, onGetDeepLink} = useDeepLink();

  return (
    <S.Container>
      <S.TitleContainer>
        <S.Title>Deep Link Usage</S.Title>
      </S.TitleContainer>
      <Typography.Paragraph>
        The deep link below enables you to run this test via a browser. It is useful as you can create dashboards to run
        particular tests interactively.
      </Typography.Paragraph>
      <FramedCodeBlock
        title="The deep link structure:"
        language="bash"
        value={deepLink}
        minHeight="80px"
        maxHeight="80px"
        actions={
          <a target="_blank" href={deepLink}>
            <S.TryItButton ghost type="primary">
              Try it
            </S.TryItButton>
          </a>
        }
      />
      <Controls onChange={onGetDeepLink} environment={environment} test={test} environmentId={environmentId} />
    </S.Container>
  );
};

export default DeepLink;
