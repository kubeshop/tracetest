import {Form} from 'antd';
import SyntaxHighlighter from 'react-syntax-highlighter';
import {arduinoLight} from 'react-syntax-highlighter/dist/cjs/styles/hljs';
import {TCollectorDataStores, TDraftDataStore} from 'types/DataStore.types';
import {CollectorConfigMap} from 'constants/CollectorConfig.constants';
import useCopy from 'hooks/useCopy';
import * as S from '../../DataStorePluginForm.styled';
import * as OtelCollectorStyles from './OpenTelemetryCollector.styled';

const OpenTelemetryCollector = () => {
  const form = Form.useFormInstance<TDraftDataStore>();
  const copy = useCopy();
  const dataStoreType = Form.useWatch('dataStoreType', form) as TCollectorDataStores;
  const example = CollectorConfigMap[dataStoreType!];

  return (
    <S.FormContainer>
      <S.FormColumn>
        <OtelCollectorStyles.SubtitleContainer>
          <S.Title>Sample Configuration</S.Title>
        </OtelCollectorStyles.SubtitleContainer>
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
  );
};

export default OpenTelemetryCollector;
