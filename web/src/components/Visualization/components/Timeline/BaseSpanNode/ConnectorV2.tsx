import * as S from '../TimelineV2.styled';

interface IProps {
  hasParent: boolean;
  leftPadding: number;
  totalChildren: number;
}

const Connector = ({hasParent, leftPadding, totalChildren}: IProps) => (
  <S.Connector height="100%" width={leftPadding + 30}>
    {hasParent && (
      <>
        <S.LineBase x1={leftPadding - 3} x2={leftPadding + 12} y1="16" y2="16" />
        <S.LineBase x1={leftPadding - 4} x2={leftPadding - 4} y1="0" y2="16.5" />
      </>
    )}

    {totalChildren > 0 ? (
      <>
        <S.LineBase x1={leftPadding + 12} x2={leftPadding + 12} y1="16" y2="32" />
        <S.RectBase x={leftPadding + 2} y="8" width="20" height="16" rx="3px" ry="3px" />
        <S.TextConnector x={leftPadding + 12} y="20" textAnchor="middle">
          {totalChildren}
        </S.TextConnector>
        <S.RectBase x={leftPadding + 2} y="8" width="20" height="16" rx="3px" ry="3px" $isTransparent />
      </>
    ) : (
      <S.CircleDot cx={leftPadding + 12} cy="16" r="3" />
    )}
  </S.Connector>
);

export default Connector;
