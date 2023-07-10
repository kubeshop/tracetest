import {Space} from 'antd';
import {useRef} from 'react';

import RuleLink from 'components/AnalyzerResult/RuleLink';
import useOnClickOutside from 'hooks/useOnClickOutside';
import {TAnalyzerError} from 'types/TestRun.types';
import * as S from './AnalyzerErrors.styled';

interface IProps {
  errors: TAnalyzerError[];
  onClose(): void;
}

const AnalyzerErrorsPopover = ({errors, onClose}: IProps) => {
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
          {errors.map(analyzerError => (
            <S.RuleContainer key={analyzerError.ruleName}>
              <S.Text strong>{analyzerError.ruleName}</S.Text>

              {analyzerError.errors.length > 1 && (
                <>
                  <div>
                    <S.Text type="secondary">{analyzerError.ruleErrorDescription}</S.Text>
                  </div>
                  <S.List>
                    {analyzerError.errors.map(error => (
                      <li key={error.value}>
                        <S.Text type="secondary">{error.value}</S.Text>
                      </li>
                    ))}
                  </S.List>
                </>
              )}

              {analyzerError.errors.length === 1 && (
                <div>
                  <S.Text type="secondary">{analyzerError.errors[0].description}</S.Text>
                </div>
              )}

              {analyzerError.ruleId && <RuleLink id={analyzerError.ruleId} isSmall />}
            </S.RuleContainer>
          ))}
        </S.Body>
      </S.Content>
    </S.Container>
  );
};

export default AnalyzerErrorsPopover;
