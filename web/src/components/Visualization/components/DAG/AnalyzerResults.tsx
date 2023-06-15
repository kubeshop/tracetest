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

              {lintError.groupedErrors.map((groupedError, index) => (
                // eslint-disable-next-line react/no-array-index-key
                <div key={index}>
                  <div>
                    <S.Text type="secondary">{groupedError.error}</S.Text>
                  </div>
                  <S.List>
                    {groupedError.values?.map(value => (
                      <li key={value}>
                        <S.Text type="secondary">{value}</S.Text>
                      </li>
                    ))}
                  </S.List>
                </div>
              ))}

              {!lintError.groupedErrors.length && (
                <S.List>
                  {lintError.errors.map((error, index) => (
                    // eslint-disable-next-line react/no-array-index-key
                    <li key={index}>
                      <S.Text type="secondary">{error}</S.Text>
                    </li>
                  ))}
                </S.List>
              )}
            </div>
          ))}
        </S.Body>
      </S.Content>
    </S.Container>
  );
};

export default AnalyzerResults;
