import {Space} from 'antd';
import {useRef} from 'react';

import useOnClickOutside from 'hooks/useOnClickOutside';
import {TLintBySpanContent} from 'services/Span.service';
import * as S from './AnalyzerResults.styled';

interface IProps {
  lintErrors: TLintBySpanContent[];
  onClose(): void;
}

const AnalyzerResults = ({lintErrors, onClose}: IProps) => {
  const ref = useRef(null);
  useOnClickOutside(ref, onClose);

  return (
    <S.Container className="nowheel nodrag" ref={ref}>
      <S.Connector />
      <S.Content>
        <Space>
          <S.ErrorIcon />
          <S.Title level={4}>Analyzer errors</S.Title>
        </Space>
        <S.Body>
          {lintErrors.map(lintError => (
            <div key={lintError.ruleName}>
              <S.Text strong>{lintError.ruleName}</S.Text>

              {lintError.errors.map((error, index) => (
                // eslint-disable-next-line react/no-array-index-key
                <div key={index}>
                  <S.Text type="secondary">- {error}</S.Text>
                </div>
              ))}
            </div>
          ))}
        </S.Body>
      </S.Content>
    </S.Container>
  );
};

export default AnalyzerResults;
