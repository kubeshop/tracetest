import styled from 'styled-components';

enum NotchColor {
  HTTP = '#B9E28C',
  DB = '#DBCBD8',
  RPC = '#9AD4D6',
  MESSAGING = '#101935',
}

const getNotchColor = (system: string) => {
  if (system.startsWith('http.method')) {
    return NotchColor.HTTP;
  }
  if (system.startsWith('db.system')) {
    return NotchColor.DB;
  }
  if (system.startsWith('rpc.system')) {
    return NotchColor.RPC;
  }
  if (system.startsWith('messaging.system')) {
    return NotchColor.MESSAGING;
  }
  return '#F49D6E';
};

export const TraceNode = styled.div<{selected: boolean}>`
  background-color: white;
  border: 2px solid ${({selected}) => (selected ? 'rgb(0, 161, 253)' : 'rgb(213, 215, 224)')};
  border-radius: 16px;
  min-width: fit-content;
  display: flex;
  width: 200px;
  height: 80px;
  justify-content: center;
  align-items: center;
`;

export const TraceNotch = styled.div<{system: string}>`
  background-color: ${({system}) => getNotchColor(system)};
  margin-top: -8px;
  position: absolute;
  top: 0px;
  padding-top: 4px;
  padding-bottom: 4px;
  padding-left: 16px;
  padding-right: 16px;
  border-radius: 16px;
  width: 70%;
  justify-content: center;
  align-items: center;
  text-align: center;
`;
