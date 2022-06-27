import {Typography} from 'antd';
import styled from 'styled-components';

import {SemanticGroupNames, SemanticGroupNamesToColor} from 'constants/SemanticGroupNames.constants';

export const Body = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  max-width: 150px;
  padding: 30px 6px 6px;
  width: 150px;
`;

export const BodyText = styled(Typography.Text).attrs({
  ellipsis: true,
})`
  font-size: 12px;
  margin: 0;
`;

export const Container = styled.div<{selected: boolean}>`
  align-items: center;
  background-color: white;
  border: 1px solid ${({selected}) => (selected ? '#48586C' : '#C9CEDB')};
  border-radius: 2px;
  display: flex;
  height: 75px;
  justify-content: center;
  max-width: 150px;
  min-width: fit-content;
  width: 150px;
`;

export const Header = styled.div<{type: SemanticGroupNames}>`
  align-items: center;
  background-color: ${({type}) => SemanticGroupNamesToColor[type]};
  border-radius: 2px 2px 0 0;
  font-weight: 700;
  margin-top: 1px;
  padding: 3px 6px;
  position: absolute;
  top: 0;
  width: 99%;
`;
