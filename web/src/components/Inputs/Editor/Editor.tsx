import {lazy, Suspense} from 'react';
import {EditorView} from '@codemirror/view';
import {BasicSetupOptions} from '@uiw/react-codemirror';
import {Extension} from '@codemirror/state';
import {SupportedEditors} from 'constants/Editor.constants';
import {Completion} from '@codemirror/autocomplete';
import {TResolveExpressionContext} from 'types/Expression.types';

const EditorMap = {
  [SupportedEditors.Expression]: lazy(() => import('./Expression')),
  [SupportedEditors.Selector]: lazy(() => import('./Selector')),
  [SupportedEditors.Interpolation]: lazy(() => import('./Interpolation')),
  [SupportedEditors.CurlCommand]: lazy(() => import('./CurlCommand')),
  [SupportedEditors.Definition]: lazy(() => import('./Definition')),
} as const;

export interface IEditorProps {
  onChange?(value: string): void;
  onFocus?(view: EditorView): void;
  value?: string;
  placeholder?: string;
  basicSetup?: BasicSetupOptions;
  editable?: boolean;
  extensions?: Extension[];
  indentWithTab?: boolean;
  autoFocus?: boolean;
  onSelectAutocompleteOption?(option: Completion): void;
  context?: TResolveExpressionContext;
  autocompleteCustomValues?: string[];
}

interface IProps extends IEditorProps {
  type: SupportedEditors;
}

const Editor = ({type, ...props}: IProps) => {
  const Component = EditorMap[type];

  return (
    <Suspense fallback={<div data-cy="editor-fallback" />}>
      <Component {...props} />
    </Suspense>
  );
};

export default Editor;
