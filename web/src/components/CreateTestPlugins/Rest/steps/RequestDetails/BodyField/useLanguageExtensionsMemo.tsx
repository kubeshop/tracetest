import {json, jsonLanguage, jsonParseLinter} from '@codemirror/lang-json';
import {xml, xmlLanguage} from '@codemirror/lang-xml';
import {LanguageSupport} from '@codemirror/language';
import {Diagnostic, linter, lintGutter} from '@codemirror/lint';
import {EditorView} from '@codemirror/view';
import {XMLValidator} from 'fast-xml-parser';
import {useMemo} from 'react';
import {BodyMode} from './useBodyMode';

export function useLanguageExtensionsMemo(bodyMode: BodyMode): any[] {
  return useMemo(() => {
    switch (bodyMode) {
      case 'xml':
        return [
          xml(),
          linter((view: EditorView): Diagnostic[] => {
            const result = XMLValidator.validate(view.state.doc.sliceString(0));
            if (result === true) return [];
            return [
              {
                actions: [],
                severity: 'error',
                source: result.err.code,
                message: result.err.msg,
                from: result.err.line,
                to: result.err.col,
              },
            ];
          }),
          new LanguageSupport(xmlLanguage),
          lintGutter({}),
        ];
      case 'json':
        return [
          json(),
          linter((view: EditorView): Diagnostic[] => jsonParseLinter()(view)),
          new LanguageSupport(jsonLanguage),
          lintGutter({}),
        ];
      case 'raw':
        return [];
      default:
        return [lintGutter({})];
    }
  }, [bodyMode]);
}
