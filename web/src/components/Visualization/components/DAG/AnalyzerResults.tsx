import {Space} from 'antd';
import {useRef} from 'react';

import useOnClickOutside from 'hooks/useOnClickOutside';
import {TLintBySpanContent} from 'models/LinterResult.model';
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
            <S.RuleContainer key={lintError.ruleName}>
              <S.Text strong>{lintError.ruleName}</S.Text>

              {lintError.errors.length > 1 && (
                <>
                  <div>
                    <S.Text type="secondary">{lintError.ruleErrorDescription}</S.Text>
                  </div>
                  <S.List>
                    {lintError.errors.map(error => (
                      <li key={error.value}>
                        <S.Text type="secondary">{error.value}</S.Text>
                      </li>
                    ))}
                  </S.List>
                </>
              )}

              {lintError.errors.length === 1 && (
                <div>
                  <S.Text type="secondary">{lintError.errors[0].description}</S.Text>
                </div>
              )}
            </S.RuleContainer>
          ))}
        </S.Body>
      </S.Content>
    </S.Container>
  );
};

export default AnalyzerResults;
