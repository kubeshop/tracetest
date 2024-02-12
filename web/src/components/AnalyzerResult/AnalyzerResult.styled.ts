import {CheckCircleFilled, CloseCircleFilled, ExclamationCircleFilled, MinusCircleFilled} from '@ant-design/icons';
import {Button, Progress, Typography} from 'antd';
import styled, {css} from 'styled-components';
import noResultsIcon from 'assets/SpanAssertionsEmptyState.svg';

export const Container = styled.div`
  padding: 24px;
  background: ${({theme}) => theme.color.white};
`;

export const Title = styled(Typography.Title)`
  && {
    margin-bottom: 8px;
    display: flex;
    align-items: center;
  }
`;

export const Description = styled(Typography.Paragraph).attrs({
  type: 'secondary',
})`
  && {
    margin-bottom: 30px;
  }
`;

export const GlobalResultWrapper = styled.div`
  display: grid;
  grid-template-columns: auto 1fr;
  margin-bottom: 28px;
  gap: 45px;
`;

export const GlobalScoreWrapper = styled.div`
  display: flex;
  flex-direction: column;
  align-items: center;
`;

export const ScoreResultWrapper = styled(GlobalScoreWrapper)`
  align-items: flex-start;
`;

export const GlobalScoreContainer = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: center;
`;

export const RuleHeader = styled.div`
  display: flex;
  flex-direction: row;
  justify-content: space-between;
`;

export const Column = styled(RuleHeader)`
  width: 95%;
`;

export const RuleBody = styled(Column)<{$resultCount: number}>`
  padding-left: 20px;
  height: ${({$resultCount}) => ($resultCount > 10 ? '100vh' : `${$resultCount * 32}px`)};
`;

export const Subtitle = styled(Typography.Title)`
  && {
    && {
      margin: 0;
      margin-bottom: 8px;
    }
  }
`;

export const ResultText = styled(Typography.Text)<{$passed: boolean}>`
  && {
    color: ${({theme, $passed}) => ($passed ? theme.color.success : theme.color.error)};
  }
`;

export const ScoreProgress = styled(Progress)`
  .ant-progress-inner {
    height: 50px !important;
    width: 50px !important;
  }

  .ant-progress-circle-trail,
  .ant-progress-circle-path {
    stroke-width: 20px;
  }
`;

interface IIconProps {
  $small: boolean;
}

const getBaseIconStyle = (color: string, {$small}: IIconProps) => css`
  color: ${color};
  cursor: pointer;
  font-size: ${$small ? '14px' : '20px'};
`;

export const PassedIcon = styled(CheckCircleFilled)<IIconProps>`
  ${({theme, $small}) => getBaseIconStyle(theme.color.success, {$small})}
`;

export const FailedIcon = styled(CloseCircleFilled)<IIconProps>`
  ${({theme, $small}) => getBaseIconStyle(theme.color.error, {$small})}
`;

export const DisableIcon = styled(MinusCircleFilled)<IIconProps>`
  ${({theme, $small}) => getBaseIconStyle(theme.color.textLight, {$small})}
`;

export const WarningIcon = styled(ExclamationCircleFilled)<IIconProps>`
  ${({theme, $small}) => getBaseIconStyle(theme.color.warningYellow, {$small})}
`;

export const SpanButton = styled(Button)<{$error?: boolean}>`
  color: ${({theme, $error}) => ($error ? theme.color.error : theme.color.success)};
  padding-left: 0;
`;

export const EmptyContainer = styled.div`
  align-items: center;
  display: flex;
  flex-direction: column;
  height: calc(100% - 70px);
  justify-content: center;
  margin-top: 50px;
`;

export const EmptyIcon = styled.img.attrs({
  src: noResultsIcon,
})`
  height: auto;
  margin-bottom: 16px;
  width: 90px;
`;

export const EmptyText = styled(Typography.Text)`
  color: ${({theme}) => theme.color.textSecondary};
`;

export const EmptyTitle = styled(Typography.Title).attrs({level: 3})``;

export const ConfigureButtonContainer = styled.div`
  margin-top: 6px;
`;

export const SwitchContainer = styled.div`
  align-items: center;
  display: flex;
  gap: 8px;
  justify-content: flex-end;
  margin-bottom: 16px;
`;

export const List = styled.ul`
  padding-inline-start: 20px;
  margin-bottom: 4px;
`;

export const RuleLinkText = styled(Typography.Text)<{$isSmall: boolean}>`
  color: ${({theme}) => theme.color.textSecondary};
  font-size: ${({theme, $isSmall}) => ($isSmall ? theme.size.xs : theme.size.sm)};
`;
