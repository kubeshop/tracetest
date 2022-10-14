import {LRLanguage, LanguageSupport} from '@codemirror/language';
import {styleTags, tags as t} from '@lezer/highlight';
import {parser} from './grammar';

export const selectorQLang = LRLanguage.define({
  parser: parser.configure({
    props: [
      styleTags({
        Identifier: t.keyword,
        String: t.string,
        Operator: t.operatorKeyword,
        Number: t.number,
        Span: t.tagName,
        ClosingBracket: t.tagName,
        Comma: t.operatorKeyword,
        PseudoSelector: t.operatorKeyword,
      }),
    ],
  }),
});

export const selectorQL = () => {
  return new LanguageSupport(selectorQLang);
};
