import React from 'react';
import {Handle, NodeProps, Position} from 'react-flow-renderer';
import * as S from './SkeletonDiagram.styled';

const SkeletonNode: React.FC<NodeProps<{}>> = ({id}) => {
  return (
    <S.SkeletonNode>
      <S.SkeletonNotch>
        <S.TextHolder />
      </S.SkeletonNotch>
      <Handle type="target" id={id} position={Position.Top} style={{top: 0, borderRadius: 0, visibility: 'hidden'}} />
      <S.TextContainer>
        <S.TextHolder />
        <S.TextHolder />
      </S.TextContainer>
      <Handle
        type="source"
        position={Position.Bottom}
        id={id}
        style={{bottom: 0, borderRadius: 0, visibility: 'hidden'}}
      />
    </S.SkeletonNode>
  );
};

export default SkeletonNode;
