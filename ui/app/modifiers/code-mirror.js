import { action } from '@ember/object';
import { bind } from '@ember/runloop';
import codemirror from 'codemirror';
import Modifier from 'ember-modifier';

import 'codemirror/addon/edit/matchbrackets';
import 'codemirror/addon/selection/active-line';
// right now we only use the ruby mode, if you use another you need to import it.
// https://codemirror.net/mode/
import 'codemirror/mode/ruby/ruby';

export default class CodeMirrorModifier extends Modifier {
  didInstall() {
    this._setup();
  }

  didUpdateArguments() {
    if (this._editor.getValue() !== this.args.named.content) {
      this._editor.setValue(this.args.named.content);
    }

    this._editor.setOption('readOnly', this.args.named.readOnly);
    this._editor.setOption('mode', this.mode);
  }

  @action
  _onChange(editor) {
    this.args.named.onUpdate(editor.getValue());
  }

  _setup() {
    if (!this.element) {
      throw new Error('CodeMirror modifier has no element');
    }

    // ARG TODO I might be able to remove those with false or empty values check with defaults on codemirror
    const editor = codemirror(this.element, {
      matchBrackets: true,
      lint: { lintOnChange: false },
      showCursorWhenSelecting: true,
      styleActiveLine: true,
      tabSize: 2,
      extraKeys: this.args.named.extraKeys || '',
      gutters: this.args.named.gutters || ['CodeMirror-lint-markers'],
      lineNumbers: this.args.named.lineNumber || true,
      mode: this.args.named.mode || 'application/json',
      readOnly: this.args.named.readOnly || false,
      theme: this.args.named.theme || 'hashi',
      value: this.args.named.content || '',
      viewportMargin: this.args.named.viewportMargin || '',
    });

    editor.on('change', bind(this, this._onChange));

    this._editor = editor;
  }
}
