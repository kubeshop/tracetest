import {Typography} from 'antd';
import styled, {css} from 'styled-components';

import {Colors} from 'constants/DAG.constants';
import {SemanticGroupNames, SemanticGroupNamesToColor} from 'constants/SemanticGroupNames.constants';

export const Badge = styled.div`
  align-items: center;
  display: flex;
  margin-top: 2px;
  color: rgba(3, 24, 73, 0.4);
`;

export const BadgeText = styled(Typography.Text)`
  font-size: 12px;
  color: inherit;
  margin-left: 4px;
`;

export const Body = styled.div`
  background-color: white;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 6px;
  width: 100%;
`;

export const BodyText = styled(Typography.Text).attrs({ellipsis: true})<{$secondary?: boolean}>`
  color: ${({$secondary}) => $secondary && 'rgba(3, 24, 73, 0.4)'};
  font-size: 12px;
`;

export const Container = styled.div<{$affected: boolean; $selected: boolean}>`
  align-items: center;
  border: ${({$selected}) => ($selected ? `1px solid ${Colors.Selected}` : `1px solid ${Colors.Default}`)};
  border-radius: 4px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  width: 180px;

  ${({$affected}) =>
    $affected &&
    css`
      border: 1px solid #031849;
      box-shadow: 2px 2px 0px #031849;
    `}
`;

export const Dot = styled.div<{$type: 'success' | 'error'}>`
  background-color: ${({$type}) => ($type === 'success' ? '#49aa19' : '#ff4d4f')};
  border-radius: 50%;
  height: 8px;
  width: 8px;
`;

export const Footer = styled.div`
  bottom: 10px;
  display: flex;
  justify-content: space-between;
  position: absolute;
  right: 6px;
  width: 22px;
`;

export const Header = styled.div<{$type: SemanticGroupNames}>`
  background-color: ${({$type}) => SemanticGroupNamesToColor[$type]};
  border-top-left-radius: 4px;
  border-top-right-radius: 4px;
  padding: 4px 24px 4px 6px;
  position: relative;
  width: 100%;
`;

export const HeaderText = styled(Typography.Paragraph).attrs({ellipsis: {rows: 2}, strong: true})`
  font-size: 12px;

  &.ant-typography {
    margin-bottom: 0;
  }
`;

export const IconContainer = styled.div<{$type: SemanticGroupNames}>`
  align-items: center;
  background-color: ${({$type}) => SemanticGroupNamesToColor[$type]};
  border: 2px solid #ffffff;
  border-radius: 50%;
  bottom: -12px;
  color: #ffffff;
  display: flex;
  font-size: 12px;
  height: 24px;
  justify-content: center;
  position: absolute;
  right: 4px;
  width: 24px;
`;
