import {SyntaxNode} from '@lezer/common/dist';
import {EditorState} from '@codemirror/state';
import {syntaxTree} from '@codemirror/language';
import {CompletionContext} from '@codemirror/autocomplete';
import {
  completeSourceAfter,
  operatorList,
  parserList,
  SourceByEditorType,
  SupportedEditors,
  Tokens,
} from 'constants/Editor.constants';
import {TSpanFlatAttribute} from 'types/Span.types';

const environmentList = [
  {
    label: 'HOST',
    type: 'variableName',
  },
  {
    label: 'PORT',
    type: 'variableName',
  },
];

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
    attributeList: TSpanFlatAttribute[] = []
  ) {
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
            ? environmentList
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
            ? environmentList
            : attributeList.map(({key}) => ({
                label: key,
                type: 'variableName',
              })),
      };
    }

    return null;
  },

  getAutocomplete(
    type: SupportedEditors.Interpolation | SupportedEditors.Expression,
    context: CompletionContext,
    attributeList: TSpanFlatAttribute[] = []
  ) {
    const {state, pos} = context;
    const tree = syntaxTree(state);
    const node = tree.resolveInner(pos, -1);

    const parserAutocomplete = this.getParserAutocomplete(node);
    if (parserAutocomplete) return parserAutocomplete;

    const operatorAutocomplete = this.getOperatorAutocomplete(node);
    if (operatorAutocomplete) return operatorAutocomplete;

    return this.getSourceAutocomplete(type, node, state, attributeList);
  },
});

export default EditorService();
