import {autocompletion} from '@codemirror/autocomplete';
import {useMemo} from 'react';
import {Attributes} from '../../constants/SpanAttribute.constants';
import {tracetest} from '../../utils/grammar';

function isNumber(text?: string) {
  return text?.toString().match(/^\d*(\.\d+)?$/);
}

export function useExpectedInputLanguage() {
  return useMemo(
    () => [
      autocompletion({
        override: [
          context => {
            const message = context.matchBefore(/\w*/);
            const attributeOptions = Object.values(Attributes).map(s => ({label: s, value: s, apply: `${s} `}));
            const durationOtions = [
              {label: `${message?.text.toString()}ms`, value: `${message?.text.toString()}ms`},
              {label: `${message?.text.toString()}s`, value: `${message?.text.toString()}s`},
            ];
            const isN = isNumber(message?.text);
            if (message?.from === 0) {
              return {from: 0, options: isN ? durationOtions : attributeOptions};
            }
            return {
              from: context.pos - 1,
              options: isN
                ? durationOtions
                : [
                    {label: '+', value: '-', apply: '- '},
                    {label: '-', value: '+', apply: '+ '},
                  ],
            };
          },
        ],
      }),
      tracetest(),
    ],
    []
  );
}
