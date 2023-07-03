import {useEffect, useRef, useState} from 'react';
import AttributeList from 'components/AttributeList';
import {OtelReference} from 'components/TestSpecForm/hooks/useGetOTELSemanticConventionAttributesInfo';
import TestRunOutput from 'models/TestRunOutput.model';
import {TSpanFlatAttribute} from 'types/Span.types';
import {TTestSpecSummary} from 'types/TestRun.types';
import {IPropsAttributeRow} from './SpanDetail';
import * as S from './SpanDetail.styled';

interface IProps {
  attributeList: TSpanFlatAttribute[];
  searchText?: string;
  semanticConventions: OtelReference;
  testSpecs?: TTestSpecSummary;
  testOutputs?: TestRunOutput[];
  onCreateTestSpec(attribute: TSpanFlatAttribute): void;
  onCreateOutput(attribute: TSpanFlatAttribute): void;
  AttributeRowComponent: React.ComponentType<IPropsAttributeRow>;
}

const Attributes = ({
  attributeList,
  searchText,
  semanticConventions,
  testSpecs,
  testOutputs,
  onCreateTestSpec,
  onCreateOutput,
  AttributeRowComponent,
}: IProps) => {
  const containerRef = useRef<HTMLDivElement>(null);
  const [topPosition, setTopPosition] = useState(0);

  useEffect(() => {
    setTopPosition(containerRef?.current?.offsetTop ?? 0);
  }, [attributeList]);

  return (
    <S.AttributesContainer $top={topPosition} ref={containerRef}>
      <AttributeList
        attributeList={attributeList}
        searchText={searchText}
        semanticConventions={semanticConventions}
        testSpecs={testSpecs}
        testOutputs={testOutputs}
        onCreateTestSpec={onCreateTestSpec}
        onCreateOutput={onCreateOutput}
        AttributeRowComponent={AttributeRowComponent}
      />
    </S.AttributesContainer>
  );
};

export default Attributes;
