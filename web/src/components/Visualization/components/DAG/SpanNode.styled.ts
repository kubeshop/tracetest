import {ExclamationCircleFilled} from '@ant-design/icons';
import {Badge, Typography} from 'antd';
import styled, {css} from 'styled-components';

import {
  SemanticGroupNames,
  SemanticGroupNamesToColor,
  SemanticGroupNamesToLightColor,
} from 'constants/SemanticGroupNames.constants';

export const BadgeCheck = styled(Badge)`
  font-size: ${({theme}) => theme.size.xs};
  line-height: 10px;
  padding: 0 2px;
  vertical-align: bottom;

  .ant-badge-status-dot {
    height: 8px;
    top: 0;
    width: 8px;
  }

  .ant-badge-status-text {
    font-size: ${({theme}) => theme.size.xs};
    margin-left: 2px;
  }
`;

export const BadgeContainer = styled.div`
  display: flex;
`;

export const BadgeType = styled(Badge)<{$hasMargin?: boolean; $type: SemanticGroupNames}>`
  &.ant-badge-not-a-wrapper:not(.ant-badge-status) {
    display: block;
    vertical-align: unset;
  }

  > sup {
    background-color: ${({$type}) => SemanticGroupNamesToLightColor[$type]};
    border-radius: 2px;
    color: ${({theme}) => theme.color.text};
    font-size: 8px;
    font-weight: 600;
    height: 12px;
    line-height: 12px;
    margin-bottom: ${({$hasMargin}) => $hasMargin && '4px'};
    text-transform: uppercase;
  }
`;

export const Body = styled.div`
  background-color: ${({theme}) => theme.color.white};
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 10px 10px;
  width: 100%;
`;

export const Container = styled.div<{$matched: boolean; $selected: boolean}>`
  align-items: center;
  background-color: ${({theme}) => theme.color.white};
  border: ${({theme, $selected}) =>
    $selected ? `2px solid ${theme.color.interactive}` : `1px solid ${theme.color.border}`};
  box-shadow: ${({$selected}) => $selected && '3px 3px 6px 0px rgba(59, 97, 246, 0.6)'};
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  overflow: hidden;
  width: 180px;

  ${({$matched, $selected}) =>
    $matched &&
    css`
      border: ${({theme}) => !$selected && `1px solid ${theme.color.text}`};
      box-shadow: ${({theme}) => !$selected && `2px 2px 0px ${theme.color.text}`};
    `}
`;

export const Footer = styled.div`
  bottom: 12px;
  position: absolute;
  flex-direction: column;
  align-items: flex-end;
  right: 10px;
  display: flex;
`;

export const Header = styled.div`
  padding: 10px 10px 4px;
  position: relative;
  width: 100%;
`;

export const HeaderText = styled(Typography.Paragraph).attrs({ellipsis: {rows: 2}, strong: true})`
  font-size: ${({theme}) => theme.size.sm};
  text-align: center;

  &.ant-typography {
    margin-bottom: 0;
  }
`;

export const Item = styled.div`
  align-items: center;
  display: flex;
  color: ${({theme}) => theme.color.text};
  font-size: ${({theme}) => theme.size.xs};
`;

export const ItemText = styled(Typography.Text)`
  color: inherit;
  margin-left: 5px;
`;

export const LintErrorIcon = styled(ExclamationCircleFilled)`
  color: ${({theme}) => theme.color.error};
  position: absolute;
  right: 10px;
  top: 5px;
`;

export const TopLine = styled.div<{$type: SemanticGroupNames}>`
  background-color: ${({$type}) => SemanticGroupNamesToColor[$type]};
  height: 7px;
  width: 100%;
`;
