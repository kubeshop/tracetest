import {LRLanguage, LanguageSupport} from '@codemirror/language';
import {styleTags, tags as t} from '@lezer/highlight';
import {parser} from './grammar';

export const interpolationQLang = LRLanguage.define({
  parser: parser.configure({
    props: [
      styleTags({
        Identifier: t.string,
        String: t.keyword,
        OpenInterpolation: t.operatorKeyword,
        CloseInterpolation: t.operatorKeyword,
        Colon: t.operatorKeyword,
        Source: t.operatorKeyword,
        Operator: t.operatorKeyword,
        Number: t.number,
      }),
    ],
  }),
});

export const interpolationQL = () => {
  return new LanguageSupport(interpolationQLang);
};
