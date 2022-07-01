import {Badge, Typography} from 'antd';
import styled, {css} from 'styled-components';

import {Colors} from 'constants/DAG.constants';
import {
  SemanticGroupNames,
  SemanticGroupNamesToColor,
  SemanticGroupNamesToLightColor,
} from 'constants/SemanticGroupNames.constants';

export const BadgeCheck = styled(Badge)`
  font-size: 10px;
  line-height: 10px;
  padding: 0 2px;
  vertical-align: bottom;

  .ant-badge-status-dot {
    height: 8px;
    top: 0;
    width: 8px;
  }

  .ant-badge-status-text {
    color: #031849;
    font-size: 10px;
    margin-left: 2px;
  }
`;

export const BadgeContainer = styled.div`
  display: flex;
`;

export const BadgeType = styled(Badge)<{$type: SemanticGroupNames}>`
  &.ant-badge-not-a-wrapper:not(.ant-badge-status) {
    display: block;
    vertical-align: unset;
  }

  > sup {
    background-color: ${({$type}) => SemanticGroupNamesToLightColor[$type]};
    border-radius: 2px;
    color: #031849;
    font-size: 8px;
    font-weight: 600;
    height: 12px;
    line-height: 12px;
    margin-bottom: 4px;
    text-transform: uppercase;
  }
`;

export const Body = styled.div`
  background-color: white;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 10px 10px;
  width: 100%;
`;

export const Container = styled.div<{$affected: boolean; $selected: boolean}>`
  align-items: center;
  background-color: #ffffff;
  border: ${({$selected}) => ($selected ? `1px solid ${Colors.Selected}` : `1px solid ${Colors.Default}`)};
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  overflow: hidden;
  width: 180px;

  ${({$affected}) =>
    $affected &&
    css`
      border: 1px solid #031849;
      box-shadow: 2px 2px 0px #031849;
    `}
`;

export const Footer = styled.div`
  bottom: 12px;
  position: absolute;
  right: 10px;
`;

export const Header = styled.div`
  padding: 10px 10px 4px;
  position: relative;
  width: 100%;
`;

export const HeaderText = styled(Typography.Paragraph).attrs({ellipsis: {rows: 2}, strong: true})`
  font-size: 12px;
  text-align: center;

  &.ant-typography {
    margin-bottom: 0;
  }
`;

export const Item = styled.div`
  align-items: center;
  display: flex;
  color: #031849;
  font-size: 10px;
`;

export const ItemText = styled(Typography.Text)`
  color: inherit;
  margin-left: 5px;
`;

export const TopLine = styled.div<{$type: SemanticGroupNames}>`
  background-color: ${({$type}) => SemanticGroupNamesToColor[$type]};
  height: 7px;
  width: 100%;
`;
