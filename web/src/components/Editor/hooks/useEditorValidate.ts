import {useCallback} from 'react';
import {SupportedEditors} from 'constants/Editor.constants';
import {expressionQLang} from '../Expression/grammar';
import {interpolationQLang} from '../Interpolation/grammar';
import {selectorQLang} from '../Selector/grammar';

const langMap = {
  [SupportedEditors.Expression]: expressionQLang,
  [SupportedEditors.Interpolation]: interpolationQLang,
  [SupportedEditors.Selector]: selectorQLang,
} as const;

const useEditorValidate = () => {
  return useCallback(
    (type: SupportedEditors.Expression | SupportedEditors.Interpolation | SupportedEditors.Selector, query: string) => {
      const lang = langMap[type];
      try {
        lang.parser.configure({strict: true}).parse(query);
        return true;
      } catch (e) {
        return false;
      }
    },
    []
  );
};

export default useEditorValidate;
