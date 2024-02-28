import {Tooltip} from 'antd';
import {ResourceType} from 'types/Resource.type';
import {abbreviateNumber} from 'utils/Common';
import * as S from './ResourceCard.styled';

const INITIAL_TIER = 1000;

interface IProps {
  num: number;
  type: ResourceType;
}

const Box = ({num, type}: IProps) => {
  if (num < INITIAL_TIER) {
    return (
      <S.Box $type={type}>
        <S.BoxTitle level={4}>{num}</S.BoxTitle>
      </S.Box>
    );
  }

  return (
    <S.Box $type={type}>
      <Tooltip title={num}>
        <S.BoxTitle level={4}>{abbreviateNumber(num)}</S.BoxTitle>
      </Tooltip>
    </S.Box>
  );
};

export default Box;
