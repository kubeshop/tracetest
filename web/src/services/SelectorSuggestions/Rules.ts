import findLast from 'lodash/findLast';

import {PseudoSelector} from 'constants/Operator.constants';
import {TSelector} from 'types/Common.types';
import {ISuggestion} from 'types/TestSpecs.types';
import {ordinalSuffixOf} from 'utils/Common';

interface IRule {
  match(
    selector: TSelector,
    matchedSpansId: string[],
    selectedSpanId: string,
    selectedSpanSelector: string,
    selectedParentSpanSelector: string
  ): ISuggestion | undefined;
}

export const allSpansRule: IRule = {
  match(selector) {
    if (!selector.query) return;
    return {query: '', title: 'All spans'};
  },
};

export const selectedSpanRule: IRule = {
  match(selector, matchedSpansId, selectedSpanId, selectedSpanSelector) {
    if (selector.query === selectedSpanSelector) return;

    return {
      query: selectedSpanSelector,
      title: 'Current span',
    };
  },
};

export const structuralPseudoClassRule: IRule = {
  match(selector, matchedSpansId, selectedSpanId) {
    if (matchedSpansId.length <= 1 || !selector.query) return;

    const foundIndex = matchedSpansId.findIndex(matchedSpanId => matchedSpanId === selectedSpanId);
    if (foundIndex === -1) return;

    const index = foundIndex + 1;
    let structuralPseudoClass = '';
    if (index === 1) {
      structuralPseudoClass = PseudoSelector.FIRST;
    } else if (index === matchedSpansId.length) {
      structuralPseudoClass = PseudoSelector.LAST;
    } else {
      structuralPseudoClass = `${PseudoSelector.NTH}(${index})`;
    }

    return {
      query: `${selector.query}${structuralPseudoClass}`,
      title: `${ordinalSuffixOf(index)} span in group`,
    };
  },
};

export function createByAttributeAllSpansRule(byAttribute: string): IRule {
  return {
    match(selector) {
      const mainFilters =
        selector.structure?.flatMap(structureItem => {
          return structureItem?.filters?.map(filter => ({...filter}));
        }) ?? [];
      const childFilters =
        selector.structure?.flatMap(structureItem => {
          return structureItem?.childSelector?.filters?.map(filter => ({...filter}));
        }) ?? [];
      const filters = [...mainFilters, ...childFilters].filter(item => item !== undefined);
      const attribute = findLast(filters, filter => filter?.property?.includes(byAttribute) ?? false);
      if (!attribute || filters.length === 1) return;

      return {
        query: `span[${byAttribute}="${attribute.value}"]`,
        title: `All ${attribute.value} spans`,
      };
    },
  };
}

export function createByAttributeDescendantsRule(byAttribute: string): IRule {
  return {
    match(selector, matchedSpansId, selectedSpanId, selectedSpanSelector, selectedParentSpanSelector) {
      if (!selectedParentSpanSelector) return;

      const filters =
        selector.structure?.flatMap(structureItem => structureItem?.filters?.map(filter => ({...filter}))) ?? [];
      const attribute = filters.find(filter => filter?.property?.includes(byAttribute) ?? false);
      if (!attribute) return;

      const childSelectors = selector.structure?.map(structureItem => 'childSelector' in (structureItem ?? {})) ?? [];
      if (childSelectors.includes(true)) return;

      return {
        query: `${selectedParentSpanSelector} span[${byAttribute}="${attribute.value}"]`,
        title: `${attribute.value} descendants`,
      };
    },
  };
}
