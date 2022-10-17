import {lazy, Suspense} from 'react';
import {BasicSetupOptions} from '@uiw/react-codemirror';
import {Extension} from '@codemirror/state';
import {SupportedEditors} from 'constants/Editor.constants';

const EditorMap = {
  [SupportedEditors.Expression]: lazy(() => import('./Expression')),
  [SupportedEditors.Selector]: lazy(() => import('./Selector')),
  [SupportedEditors.Interpolation]: lazy(() => import('./Interpolation')),
  [SupportedEditors.CurlCommand]: lazy(() => import('./CurlCommand')),
} as const;

export interface IEditorProps {
  onChange?(value: string): void;
  value?: string;
  placeholder?: string;
  basicSetup?: BasicSetupOptions;
  editable?: boolean;
  extensions?: Extension[];
  indentWithTab?: boolean;
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
