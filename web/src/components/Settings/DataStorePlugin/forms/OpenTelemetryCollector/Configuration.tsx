import {Form} from 'antd';
import SyntaxHighlighter from 'react-syntax-highlighter';
import {arduinoLight} from 'react-syntax-highlighter/dist/cjs/styles/hljs';
import useCopy from 'hooks/useCopy';
import {CollectorConfigMap} from 'constants/CollectorConfig.constants';
import {TCollectorDataStores, TDraftDataStore} from 'types/DataStore.types';
import DocsBanner from 'components/DocsBanner/DocsBanner';
import * as S from '../../DataStorePluginForm.styled';
import * as OtelCollectorStyles from './OpenTelemetryCollector.styled';

const Configuration = () => {
  const copy = useCopy();
  const form = Form.useFormInstance<TDraftDataStore>();
  const dataStoreType = Form.useWatch('dataStoreType', form) as TCollectorDataStores;
  const example = CollectorConfigMap[dataStoreType!];

  return (
    <>
      <OtelCollectorStyles.SubtitleContainer>
        <OtelCollectorStyles.Title>OpenTelemetry Collector Configuration</OtelCollectorStyles.Title>
        <OtelCollectorStyles.Description>
          The OpenTelemetry Collector configuration below is a sample. Your config file layout should look the same.
          Make sure to use your own API keys or tokens as explained. Copy the config sample below, paste it into your
          own OpenTelemetry Collector config and apply it.
        </OtelCollectorStyles.Description>
      </OtelCollectorStyles.SubtitleContainer>
      <S.FormContainer>
        <S.FormColumn>
          <OtelCollectorStyles.CodeContainer data-cy="file-viewer-code-container">
            <OtelCollectorStyles.CopyIconContainer onClick={() => copy(example)}>
              <OtelCollectorStyles.CopyIcon />
            </OtelCollectorStyles.CopyIconContainer>
            <SyntaxHighlighter language="yaml" style={arduinoLight}>
              {example}
            </SyntaxHighlighter>
          </OtelCollectorStyles.CodeContainer>
        </S.FormColumn>
      </S.FormContainer>
      <DocsBanner>
        Need more information about setting up OpenTelemetry Collector? <a href="">Go to our docs</a>
      </DocsBanner>
    </>
  );
};

export default Configuration;
