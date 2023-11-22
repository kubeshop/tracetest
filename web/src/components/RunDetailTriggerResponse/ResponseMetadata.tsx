import {useTestRun} from 'providers/TestRun/TestRun.provider';
import {useMemo} from 'react';
import * as S from './RunDetailTriggerResponse.styled';
import KeyValueRow from '../KeyValueRow';

const ResponseMetadata = () => {
  const {
    run: {metadata},
  } = useTestRun();

  const entries = useMemo(() => Object.entries(metadata).filter(([, value]) => !!value), [metadata]);

  if (!entries.length) {
    return (
      <S.EmptyContainer>
        <S.EmptyIcon />
        <S.EmptyTitle>There are no metadata entries used in this test</S.EmptyTitle>
      </S.EmptyContainer>
    );
  }

  return (
    <S.ResponseVarsContainer>
      {entries.map(([key, value]) => (
        <KeyValueRow key={key} keyName={key} value={value} />
      ))}
    </S.ResponseVarsContainer>
  );
};

export default ResponseMetadata;
