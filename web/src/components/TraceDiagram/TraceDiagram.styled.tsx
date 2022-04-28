import {Typography} from 'antd';
import styled from 'styled-components';
import {SemanticGroupNames} from '../../constants/SemanticGroupNames.constants';

enum NotchColor {
  HTTP = '#B9E28C',
  DB = '#DBCBD8',
  RPC = '#9AD4D6',
  MESSAGING = '#101935',
  DEFAULT = '#F49D6E',
}

const getNotchColor = (spanType: SemanticGroupNames) => {
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
    case SemanticGroupNames.Messaging: {
      return 'white';
    }
    default: {
      return 'inherit';
    }
  }
};

export const TraceNode = styled.div<{selected: boolean}>`
  background-color: white;
  border: 2px solid ${({selected}) => (selected ? 'rgb(0, 161, 253)' : 'rgb(213, 215, 224)')};
  border-radius: 4px;
  min-width: fit-content;
  display: flex;
  width: 200px;
  height: 90px;
  justify-content: center;
  align-items: center;
`;

export const TraceNotch = styled.div<{spanType: SemanticGroupNames}>`
  background-color: ${({spanType}) => getNotchColor(spanType)};
  margin-top: -8px;
  position: absolute;
  top: 0px;
  padding-top: 4px;
  padding-bottom: 4px;
  padding-left: 16px;
  padding-right: 16px;
  border-radius: 4px;
  width: 70%;
  justify-content: center;
  align-items: center;
  text-align: center;

  span {
    color: ${({spanType}) => getTextColor(spanType)};
  }
`;

export const NameText = styled(Typography.Text).attrs({
  ellipsis: true,
})`
  margin: 0;
  margin-bottom: 5px;
`;

export const TextContainer = styled.div`
  padding: 14px;
  padding-top: 35px;
  max-width: 180px;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
`;

export const Container = styled.div`
  position: relative;
`;
