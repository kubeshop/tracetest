import {useCallback} from 'react';
import {EditorView} from '@codemirror/view';
import {syntaxTree} from '@codemirror/language';
import {Tokens} from '../../../constants/Editor.constants';

const useTooltip = () => {
  return useCallback(({state}: EditorView, pos: number, side: -1 | 1) => {
    const tree = syntaxTree(state);
    const node = tree.resolveInner(pos, -1);

    if ((node.from === pos && side < 0) || (node.from === pos && side > 0)) return null;

    const parentNode = node.parent;

    if (parentNode && parentNode.name === Tokens.OutsideInput) {
      const source = parentNode.firstChild || {from: 0, to: 0};
      const identifier = parentNode.lastChild || {from: 0, to: 0};
      const [sourceText] = state.doc.sliceString(source.from, source.to).split(':');
      const identifierText = state.doc.sliceString(identifier.from, identifier.to);

      // TODO: display the value from the selected source and value
      // eslint-disable-next-line no-console
      console.log('@@', sourceText);

      return {
        pos: parentNode.from,
        end: parentNode.to,
        above: true,
        create() {
          const dom = document.createElement('div');
          dom.textContent = identifierText;
          return {dom};
        },
      };
    }

    return null;
  }, []);
};

export default useTooltip;
