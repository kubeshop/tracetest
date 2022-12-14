import {message} from 'antd';
import SyntaxHighlighter from 'react-syntax-highlighter';
import {arduinoLight} from 'react-syntax-highlighter/dist/cjs/styles/hljs';
import {OtelCollectorConfigSample} from 'constants/DataStore.constants';
import * as S from '../../DataStorePluginForm.styled';
import * as OtelCollectorStyles from './OpenTelemetryCollector.styled';

const OpenTelemetryCollector = () => {
  const onCopy = () => {
    message.success('Content copied to the clipboard');
    navigator.clipboard.writeText(OtelCollectorConfigSample);
  };

  return (
    <S.FormContainer>
      <S.FormColumn>
        <OtelCollectorStyles.SubtitleContainer>
          <S.Title>Sample Configuration</S.Title>
        </OtelCollectorStyles.SubtitleContainer>
        <OtelCollectorStyles.CodeContainer data-cy="file-viewer-code-container">
          <OtelCollectorStyles.CopyIconContainer onClick={onCopy}>
            <OtelCollectorStyles.CopyIcon />
          </OtelCollectorStyles.CopyIconContainer>
          <SyntaxHighlighter language="yaml" style={arduinoLight}>
            {OtelCollectorConfigSample}
          </SyntaxHighlighter>
        </OtelCollectorStyles.CodeContainer>
      </S.FormColumn>
    </S.FormContainer>
  );
};

export default OpenTelemetryCollector;
