import {autocompletion} from '@codemirror/autocomplete';
import {syntaxTree} from '@codemirror/language';

import {Attributes} from 'constants/SpanAttribute.constants';
import {useMemo} from 'react';
import {selectorQL} from 'components/Inputs/Editor/Selector/grammar';

function isNumber(text?: string) {
  const matches = text?.toString().match(/^\d*(\.\d+)?$/);
  return Boolean(matches?.[0]);
}

export function useExpectedInputLanguage() {
  return useMemo(
    () => [
      autocompletion({
        override: [
          context => {
            const {state, pos} = context;
            const word = context.matchBefore(/\w*/);

            const tree = syntaxTree(state);
            const nodeBefore = tree.resolveInner(pos, -1);
            const parentNode = nodeBefore?.parent ?? {from: 0, to: 0};
            const identifierText = state.doc.sliceString(parentNode.from, parentNode.to);
            const isN = isNumber(identifierText);

            const attributeOptions = Object.values(Attributes).map(s => ({label: `attr:${s}`, apply: `attr:${s} `}));
            const durationOptions = [{label: `${word?.text.toString()}ms`}, {label: `${word?.text.toString()}s`}];
            const operatorOptions = [
              {label: '+', apply: '+ '},
              {label: '-', apply: '- '},
            ];

            if (word?.from === 0) {
              return {from: 0, options: isN ? durationOptions : attributeOptions};
            }

            return {
              from: word?.from ?? 0,
              options: isN ? durationOptions : operatorOptions,
            };
          },
        ],
      }),
      selectorQL(),
    ],
    []
  );
}
