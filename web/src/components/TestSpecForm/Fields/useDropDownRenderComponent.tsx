import parse from 'html-react-parser';
import MarkdownIt from 'markdown-it';
import React, {useMemo} from 'react';
import SpanAttributeService from 'services/SpanAttribute.service';
import {OtelReference} from '../hooks/useGetOTELSemanticConventionAttributesInfo';
import * as S from './AttributeField.styled';

export function useDropDownRenderComponent(
  reference: OtelReference,
  hoveredKey?: string
): (menu: React.ReactElement) => React.ReactElement {
  const otelReferenceModel = SpanAttributeService.referencePicker(reference, hoveredKey || '');
  const description = useMemo(() => parse(MarkdownIt().render(otelReferenceModel.description)), [otelReferenceModel]);
  return menu => {
    return hoveredKey ? (
      <S.Container>
        <S.Sides>{menu}</S.Sides>
        <S.Sides>
          <S.Content>
            <S.Title>{hoveredKey}</S.Title>
            {otelReferenceModel.description !== '' ? description : 'There is no description for this attribute yet.'}
            <S.TagsContainer>
              {otelReferenceModel.tags.map(tag => (
                <S.Tag key={tag}>{tag}</S.Tag>
              ))}
            </S.TagsContainer>
          </S.Content>
        </S.Sides>
      </S.Container>
    ) : (
      menu
    );
  };
}
