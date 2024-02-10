import {BaseLeftPaddingV2} from 'constants/Timeline.constants';
import * as S from '../TimelineV2.styled';

interface IProps {
  hasParent: boolean;
  id: string;
  isCollapsed: boolean;
  nodeDepth: number;
  onCollapse(id: string): void;
  totalChildren: number;
}

const Connector = ({hasParent, id, isCollapsed, nodeDepth, onCollapse, totalChildren}: IProps) => {
  const leftPadding = nodeDepth * BaseLeftPaddingV2;

  return (
    <S.Connector height="100%" width={leftPadding + 30}>
      {hasParent && (
        <>
          <S.LineBase x1={leftPadding - 3} x2={leftPadding + 12} y1="16" y2="16" />
          <S.LineBase x1={leftPadding - 4} x2={leftPadding - 4} y1="0" y2="16.5" />
        </>
      )}

      {totalChildren > 0 ? (
        <>
          {!isCollapsed && <S.LineBase x1={leftPadding + 12} x2={leftPadding + 12} y1="16" y2="32" />}
          <S.RectBase x={leftPadding + 2} y="8" width="20" height="16" rx="3px" ry="3px" $isActive={isCollapsed} />
          <S.TextConnector x={leftPadding + 12} y="20" textAnchor="middle" $isActive={isCollapsed}>
            {totalChildren}
          </S.TextConnector>
          <S.RectBaseTransparent
            x={leftPadding + 2}
            y="8"
            width="20"
            height="16"
            rx="3px"
            ry="3px"
            onClick={event => {
              event.stopPropagation();
              onCollapse(id);
            }}
          />
        </>
      ) : (
        <S.CircleDot cx={leftPadding + 12} cy="16" r="3" />
      )}

      {new Array(nodeDepth).fill(0).map((_, index) => {
        return <S.LineBase x1={index * BaseLeftPaddingV2 + 12} x2={index * BaseLeftPaddingV2 + 12} y1="0" y2="32" />;
      })}
    </S.Connector>
  );
};

export default Connector;
