// @ts-ignore
import {fileTests} from '@lezer/generator/dist/test.cjs';
import * as fs from 'fs';
import * as path from 'path';
import {interpolationQLang} from '../index';

describe('interpolationQLanguage', () => {
  describe('cases.txt', () => {
    const testList = fileTests(fs.readFileSync(path.join(__dirname, './cases.txt'), 'utf8'), 'cases.txt');

    testList.forEach(({name, run}: {name: string; run: Function}) => {
      it(name, () => run(interpolationQLang.parser));
    });
  });
});
