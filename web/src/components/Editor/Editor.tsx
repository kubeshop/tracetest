import {noop} from 'lodash';
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
}

interface IProps extends IEditorProps {
  type: SupportedEditors;
}

const Editor = ({type, onChange = noop, value = '', placeholder, basicSetup = {}, editable = true, extensions = []}: IProps) => {
  const Component = EditorMap[type];

  return (
    <Suspense fallback={<div />}>
      <Component
        onChange={onChange}
        value={value}
        placeholder={placeholder}
        basicSetup={basicSetup}
        editable={editable}
        extensions={extensions}
      />
    </Suspense>
  );
};

export default Editor;
