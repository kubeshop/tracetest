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

export const getNotchColor = (spanType: SemanticGroupNames) => {
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
  border: 1px solid ${({selected}) => (selected ? '#48586C' : '#E2E4E6')};
  border-radius: 2px;
  min-width: fit-content;
  display: flex;
  width: 150px;
  max-width: 150px;
  height: 90px;
  justify-content: center;
  align-items: center;
`;

export const TraceNotch = styled.div<{spanType: SemanticGroupNames}>`
  background-color: ${({spanType}) => getNotchColor(spanType)};
  position: absolute;
  top: 0px;
  margin-top: 1px;
  padding: 3px 6px;
  border-radius: 2px;
  border-bottom-left-radius: 0px;
  border-bottom-right-radius: 0px;
  width: 99%;
  font-weight: 700;
  align-items: center;

  span {
    color: ${({spanType}) => getTextColor(spanType)};
  }
`;

export const NameText = styled(Typography.Text).attrs({
  ellipsis: true,
})`
  margin: 0;
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
