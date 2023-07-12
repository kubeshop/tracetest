import {NodeHeight} from 'constants/Timeline.constants';
import * as S from '../Timeline.styled';

interface IProps {
  distance: number;
  leftPadding: number;
}

const Connector = ({distance, leftPadding}: IProps) => {
  const connectorX = distance * NodeHeight - 24;

  return (
    <>
      <S.LineConnector strokeDasharray="3" x1={leftPadding - 6} x2={leftPadding - 6} y1={17} y2={-connectorX} />
      <S.LineConnector strokeDasharray="3" x1={leftPadding - 6} x2={leftPadding + 6} y1={17} y2={17} />
    </>
  );
};

export default Connector;
