import {Diagnostic} from '@codemirror/lint';
import {EditorView} from '@codemirror/view';
import {Text} from '@codemirror/state';

const CUSTOM_SYNTAX_ERRORS_WHITE_LIST = [/Unexpected token \'\$\'/, /\${env:([\w-]+)}/];

/// Calls
/// [`JSON.parse`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON/parse)
/// on the document and, if that throws an error, reports it as a
/// single diagnostic.
export const jsonParseLinter =
  () =>
  (view: EditorView): Diagnostic[] => {
    try {
      JSON.parse(view.state.doc.toString());
    } catch (e) {
      if (!(e instanceof SyntaxError)) throw e;

      const whiteListTest = CUSTOM_SYNTAX_ERRORS_WHITE_LIST.map(regex => regex.test(e.message));
      if (whiteListTest.every(t => t)) return [];

      const pos = getErrorPosition(e, view.state.doc);
      return [
        {
          from: pos,
          message: e.message,
          severity: 'error',
          to: pos,
        },
      ];
    }
    return [];
  };

function getErrorPosition(error: SyntaxError, doc: Text): number {
  let m;
  // eslint-disable-next-line no-cond-assign
  if ((m = error.message.match(/at position (\d+)/))) return Math.min(Number(m[1]), doc.length);
  // eslint-disable-next-line no-cond-assign
  if ((m = error.message.match(/at line (\d+) column (\d+)/)))
    return Math.min(doc.line(Number(m[1])).from + Number(m[2]) - 1, doc.length);
  return 0;
}
