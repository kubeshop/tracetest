import {Diagnostic} from '@codemirror/lint';
import {EditorState, Text} from '@codemirror/state';
import {syntaxTree} from '@codemirror/language';
import {EditorView} from '@codemirror/view';
import {Completion, CompletionContext, CompletionResult} from '@codemirror/autocomplete';
import {SyntaxNode} from '@lezer/common/dist';
import {
  completeSourceAfter,
  operatorList,
  parserList,
  SourceByEditorType,
  SupportedEditors,
  Tokens,
} from 'constants/Editor.constants';
import {TSpanFlatAttribute} from 'types/Span.types';
import {expressionQLang} from 'components/Editor/Expression/grammar';
import {interpolationQLang} from 'components/Editor/Interpolation/grammar';
import {selectorQLang} from 'components/Editor/Selector/grammar';
import {IKeyValue} from 'constants/Test.constants';
import {noop} from 'lodash';
import AssertionService from './Assertion.service';

const langMap = {
  [SupportedEditors.Expression]: expressionQLang,
  [SupportedEditors.Interpolation]: interpolationQLang,
  [SupportedEditors.Selector]: selectorQLang,
} as const;

interface IAutoCompleteProps {
  type: SupportedEditors.Interpolation | SupportedEditors.Expression;
  context: CompletionContext;
  attributeList?: TSpanFlatAttribute[];
  envEntryList?: IKeyValue[];
  customValueList?: string[];
  onSelect?(option: Completion): void;
}

const CUSTOM_SYNTAX_ERRORS_WHITE_LIST = [/Unexpected token \'\$\'/, /\${env:([\w-]+)}/];

function getJsonErrorPosition(error: SyntaxError, doc: Text): number {
  let m;
  // eslint-disable-next-line no-cond-assign
  if ((m = error.message.match(/at position (\d+)/))) return Math.min(Number(m[1]), doc.length);
  // eslint-disable-next-line no-cond-assign
  if ((m = error.message.match(/at line (\d+) column (\d+)/)))
    return Math.min(doc.line(Number(m[1])).from + Number(m[2]) - 1, doc.length);
  return 0;
}

const EditorService = () => ({
  getOperatorAutocomplete(node: SyntaxNode) {
    const operatorNodeOne = node.lastChild?.lastChild;
    if (operatorNodeOne?.name === Tokens.ComposedValue) {
      return {
        from: node.to,
        options: operatorList,
      };
    }

    const operatorNodeTwo = node.firstChild?.lastChild;
    if (operatorNodeTwo?.name === Tokens.ComposedValue) {
      return {
        from: node.to,
        options: operatorList,
      };
    }

    return null;
  },

  getParserAutocomplete(node: SyntaxNode) {
    if (node.name === Tokens.Pipe) {
      return {
        from: node.to,
        options: parserList,
      };
    }

    return null;
  },

  getSourceAutocomplete(
    type: SupportedEditors.Interpolation | SupportedEditors.Expression,
    node: SyntaxNode,
    state: EditorState,
    environmentList: IKeyValue[] = [],
    attributeList: TSpanFlatAttribute[] = [],
    customValueList: string[] = [],
    onSelect: (option: Completion) => void = noop
  ): CompletionResult | null {
    if (node.name === Tokens.OpenInterpolation) {
      const sourceOptionList = SourceByEditorType[type];

      return {
        from: node.to,
        options: sourceOptionList,
      };
    }

    if (completeSourceAfter.includes(node.name)) {
      const {from, to} = node;
      const [sourceText] = state.doc.sliceString(from, to).split(':');

      return {
        from: node.to,
        options:
          sourceText === 'env'
            ? environmentList.map(({key}) => ({
                label: key,
                type: 'variableName',
              }))
            : attributeList.map(({key}) => ({
                label: key,
                type: 'variableName',
              })),
      };
    }

    if (node.prevSibling?.name === Tokens.Source) {
      const {from, to} = node.prevSibling || {from: 0, to: 0};
      const [sourceText] = state.doc.sliceString(from, to).split(':');

      return {
        to: node.to,
        from: node.from,
        options:
          sourceText === 'env'
            ? environmentList.map(({key}) => ({
                label: key,
                type: 'variableName',
              }))
            : attributeList.map(({key}) => ({
                label: key,
                type: 'variableName',
              })),
      };
    }

    return {
      from: 0,
      options: customValueList.length
        ? customValueList.map(value => ({
            label: value,
            type: 'variableName',
            apply: AssertionService.extractExpectedString(value),
          }))
        : attributeList.map(({key}) => ({
            label: `attr:${key}`,
            type: 'variableName',
            apply(view: EditorView, completion: Completion, from: number, to: number) {
              onSelect(completion);

              return view.dispatch({
                changes: {from, to, insert: completion.label},
              });
            },
          })),
    };
  },

  getAutocomplete({
    type,
    context,
    attributeList = [],
    envEntryList = [],
    customValueList = [],
    onSelect = noop,
  }: IAutoCompleteProps): CompletionResult | null {
    const {state, pos} = context;
    const tree = syntaxTree(state);
    const node = tree.resolveInner(pos, -1);

    const parserAutocomplete = this.getParserAutocomplete(node);
    if (parserAutocomplete) return parserAutocomplete;

    const operatorAutocomplete = this.getOperatorAutocomplete(node);
    if (operatorAutocomplete) return operatorAutocomplete;

    return this.getSourceAutocomplete(type, node, state, envEntryList, attributeList, customValueList, onSelect);
  },

  getIsQueryValid(
    type: SupportedEditors.Expression | SupportedEditors.Interpolation | SupportedEditors.Selector,
    query: string
  ) {
    const lang = langMap[type];

    try {
      lang.parser.configure({strict: true}).parse(query);
      return true;
    } catch (e) {
      return false;
    }
  },

  /// Calls
  /// [`JSON.parse`](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/JSON/parse)
  /// on the document and, if that throws an error, reports it as a
  /// single diagnostic.
  jsonParseLinter() {
    return (view: EditorView): Diagnostic[] => {
      try {
        JSON.parse(view.state.doc.toString());
      } catch (e) {
        if (!(e instanceof SyntaxError)) throw e;

        const message = e.message;
        const whiteListTest = CUSTOM_SYNTAX_ERRORS_WHITE_LIST.map(regex => regex.test(message));
        if (whiteListTest.every(t => t)) return [];

        const pos = getJsonErrorPosition(e, view.state.doc);
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
  },
});

export default EditorService();
