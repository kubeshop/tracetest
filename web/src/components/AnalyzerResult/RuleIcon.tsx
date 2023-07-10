import {Tooltip} from 'antd';
import {LinterRuleErrorLevel} from 'models/Linter.model';
import * as S from './AnalyzerResult.styled';

interface IProps {
  level: LinterRuleErrorLevel;
  passed: boolean;
}

const ruleIconMap = {
  [LinterRuleErrorLevel.ERROR]: S.FailedIcon,
  [LinterRuleErrorLevel.WARNING]: S.WarningIcon,
  [LinterRuleErrorLevel.DISABLED]: S.DisableIcon,
} as const;

const ruleTooltipMap = {
  [LinterRuleErrorLevel.ERROR]:
    'Based on the analyzer configuration. Errored rules can fail the test run based on score or passed status',
  [LinterRuleErrorLevel.WARNING]: 'Warning rules are executed but do not fail the test run and do not impact the score',
  [LinterRuleErrorLevel.DISABLED]: 'Disabled rules are not executed and do not impact the score',
} as const;

const RuleIcon = ({passed, level}: IProps) => {
  if (passed && level === LinterRuleErrorLevel.ERROR) {
    return (
      <Tooltip title="Passing rules are considered as part of the score and the overall test result">
        <S.PassedIcon $small />
      </Tooltip>
    );
  }

  const Icon = ruleIconMap[level];

  return (
    <Tooltip title={ruleTooltipMap[level]}>
      <Icon $small />
    </Tooltip>
  );
};

export default RuleIcon;
