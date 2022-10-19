import {PseudoSelector} from 'constants/Operator.constants';
import {TSelector} from 'types/Common.types';
import {ISuggestion} from 'types/TestSpecs.types';

interface IRule {
  match(selector: TSelector): ISuggestion | undefined;
}

interface IStructuralPseudoClassTitle extends Record<PseudoSelector, (value: string) => string> {}

const structuralPseudoClassTitle: IStructuralPseudoClassTitle = {
  [PseudoSelector.FIRST]: (value: string) => `First ${value} span`,
  [PseudoSelector.LAST]: (value: string) => `Last ${value} span`,
  [PseudoSelector.NTH]: (value: string) => `nth_child ${value} span`,
  [PseudoSelector.ALL]: (value: string) => `All ${value} spans`,
};

export const allSpansRule: IRule = {
  match(selector) {
    if (!selector.query) return;
    return {query: '', title: 'All spans'};
  },
};

export function createByAttributeRule(
  byAttribute: string,
  structuralPseudoClass: PseudoSelector = PseudoSelector.ALL
): IRule {
  return {
    match(selector) {
      const filters =
        selector.structure?.flatMap(structureItem => structureItem?.filters?.map(filter => ({...filter}))) ?? [];
      const attribute = filters.find(filter => filter?.property?.includes(byAttribute));

      if (!attribute || filters.length === 1) return;

      return {
        query: `span[${byAttribute}="${attribute.value}"]${structuralPseudoClass}`,
        title: structuralPseudoClassTitle[structuralPseudoClass](attribute.value),
      };
    },
  };
}
