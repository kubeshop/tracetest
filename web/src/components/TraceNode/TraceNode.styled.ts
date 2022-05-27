import {Typography} from 'antd';
import styled from 'styled-components';
import {SemanticGroupNames} from '../../constants/SemanticGroupNames.constants';

enum NotchColor {
  HTTP = '#C1E095',
  DB = '#EFDBFF',
  RPC = '#9AD4D6',
  MESSAGING = '#BFBFBF',
  DEFAULT = '#FFBB96',
}

export const getNotchColor = (spanType: SemanticGroupNames): string => {
  switch (spanType) {
    case SemanticGroupNames.Http: {
      return NotchColor.HTTP;
    }

    case SemanticGroupNames.Database: {
      return NotchColor.DB;
    }

    case SemanticGroupNames.Rpc: {
      return NotchColor.RPC;
    }

    case SemanticGroupNames.Messaging: {
      return NotchColor.MESSAGING;
    }
    default: {
      return NotchColor.DEFAULT;
    }
  }
};

const getTextColor = (spanType: SemanticGroupNames) => {
  switch (spanType) {
    default: {
      return 'inherit';
    }
  }
};

export const TraceNode = styled.div<{selected: boolean}>`
  background-color: white;
  border: 1px solid ${({selected}) => (selected ? '#48586C' : '#C9CEDB')};
  border-radius: 2px;
  min-width: fit-content;
  display: flex;
  width: 150px;
  max-width: 150px;
  height: 75px;
  justify-content: center;
  align-items: center;
`;

export const TraceNotch = styled.div<{$spanType: SemanticGroupNames}>`
  background-color: ${({$spanType}) => getNotchColor($spanType)};
  position: absolute;
  top: 0;
  margin-top: 1px;
  padding: 3px 6px;
  border-radius: 2px 2px 0 0;
  width: 99%;
  font-weight: 700;
  align-items: center;

  span {
    color: ${({$spanType}) => getTextColor($spanType)};
  }
`;

export const NameText = styled(Typography.Text).attrs({
  ellipsis: true,
})`
  margin: 0;
  font-size: 12px;
`;

export const TextContainer = styled.div`
  padding: 6px;
  padding-top: 30px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  width: 150px;
  max-width: 150px;
`;
