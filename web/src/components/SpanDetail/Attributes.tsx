import {useEffect, useRef, useState} from 'react';

import AttributeList from 'components/AttributeList';
import {OtelReference} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import {TResultAssertions} from 'types/Assertion.types';
import {TSpanFlatAttribute} from 'types/Span.types';
import * as S from './SpanDetail.styled';

interface IProps {
  assertions?: TResultAssertions;
  attributeList: TSpanFlatAttribute[];
  searchText?: string;
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  semanticConventions: OtelReference;
}

const Attributes = ({assertions, attributeList, onCreateTestSpec, searchText, semanticConventions}: IProps) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [topPosition, setTopPosition] = useState(0);

  useEffect(() => {
    setTopPosition(containerRef?.current?.offsetTop ?? 0);
  }, [attributeList]);

  return (
    <S.AttributesContainer $top={topPosition} ref={containerRef}>
      <AttributeList
        assertions={assertions}
        attributeList={attributeList}
        onCreateTestSpec={onCreateTestSpec}
        searchText={searchText}
        semanticConventions={semanticConventions}
      />
    </S.AttributesContainer>
  );
};

export default Attributes;
