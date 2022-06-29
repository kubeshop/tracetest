import {Typography} from 'antd';
import styled from 'styled-components';

import {Colors} from 'constants/DAG.constants';
import {
  SemanticGroupNames,
  SemanticGroupNamesToColor,
  SemanticGroupNamesToDarkColor,
} from 'constants/SemanticGroupNames.constants';

export const Badge = styled.div`
  align-items: center;
  border-radius: 6px;
  border: 1px solid #dce0e0;
  display: flex;
  margin-top: 4px;
  padding: 2px 4px;
  width: fit-content;
`;

export const BadgeText = styled(Typography.Text)`
  font-size: 10px;
`;

export const Body = styled.div`
  background-color: white;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 4px 8px;
  width: 100%;
`;

export const BodyText = styled(Typography.Text).attrs({ellipsis: true})`
  font-size: 12px;
`;

export const Container = styled.div<{$selected: boolean; $type: SemanticGroupNames}>`
  align-items: center;
  background-color: ${({$type}) => SemanticGroupNamesToColor[$type]};
  border: ${({$selected}) => $selected && `2px solid ${Colors.Selected}`};
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  padding: 0 3px 3px;
  width: 180px;
`;

export const Footer = styled.div`
  display: flex;
  justify-content: flex-end;
  padding: 4px 0 0;
  width: 100%;
`;

export const Header = styled.div`
  align-items: center;
  display: flex;
  height: 44px;
  padding: 4px 8px;
  width: 100%;
`;

export const HeaderText = styled(Typography.Paragraph)<{$type: SemanticGroupNames}>`
  color: ${({$type}) => SemanticGroupNamesToDarkColor[$type]};
  flex: 1;
  font-size: 12px;

  &.ant-typography {
    margin-bottom: 0;
  }
`;

export const Logo = styled.div<{$type: SemanticGroupNames}>`
  align-items: center;
  color: ${({$type}) => SemanticGroupNamesToDarkColor[$type]};
  display: flex;
  font-size: 18px;
  justify-content: left;
  width: 30px;
`;

const CenteredDiv = styled.div`
  align-items: center;
  display: flex;
  justify-content: center;
`;

export const DotContainer = styled(CenteredDiv)`
  border: 2px solid rgba(0, 0, 0, 0.1);
  border-radius: 50%;
  height: 20px;
  margin-right: 2px;
  width: 20px;
`;

export const DotContent1 = styled(CenteredDiv)`
  background-color: #ffffff;
  border-radius: 50%;
  height: 16px;
  width: 16px;
`;

export const DotContent2 = styled(CenteredDiv)<{$type: 'success' | 'error'}>`
  border: 1px solid ${({$type}) => ($type === 'success' ? '#49aa19' : '#ff4d4f')};
  border-radius: 50%;
  height: 14px;
  width: 14px;
`;

export const Dot = styled.div<{$type: 'success' | 'error'}>`
  background-color: ${({$type}) => ($type === 'success' ? '#49aa19' : '#ff4d4f')};
  border-radius: 50%;
  height: 6px;
  width: 6px;
`;
