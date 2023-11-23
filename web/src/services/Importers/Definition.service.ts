import {load} from 'js-yaml';
import {IDefinitionValues, IImportService} from 'types/Test.types';
import {TriggerTypeToPlugin} from 'constants/Plugins.constants';
import Test, {TRawTestResource} from 'models/Test.model';
import TestService from '../Test.service';
import {isJson} from '../../utils/Common';

interface IDefinitionImportService extends IImportService {
  validate(definition: string): string[];
  load(raw: string): TRawTestResource;
}

const DefinitionImportService = (): IDefinitionImportService => ({
  async getRequest(values) {
    const {definition} = values as IDefinitionValues;

    const test = Test.FromDefinition(definition);
    return TestService.getInitialValues(test);
  },

  async validateDraft(draft) {
    const {definition} = draft as IDefinitionValues;

    return !this.validate(definition).length;
  },

  getPlugin(draft) {
    const {definition} = draft as IDefinitionValues;
    const {spec = {}}: TRawTestResource = this.load(definition);

    const triggerType = spec.trigger?.type;
    return TriggerTypeToPlugin[triggerType!];
  },

  load(raw) {
    if (isJson(raw)) return JSON.parse(raw) as TRawTestResource;

    return load(raw) as TRawTestResource;
  },

  validate(raw) {
    const {type, spec = {}}: TRawTestResource = this.load(raw);
    const errors = [];

    if (type !== 'Test') {
      errors.push(`Invalid type: ${type}`);
    }

    if (!spec.name) {
      errors.push(`Missing Name`);
    }

    if (!spec.trigger) {
      errors.push(`Missing trigger`);
    }

    if (!spec.trigger?.type) {
      errors.push('Missing trigger type');
    }

    return errors;
  },
});

export default DefinitionImportService();
