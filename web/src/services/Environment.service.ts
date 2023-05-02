import Environment, {TRawEnvironment} from 'models/Environment.model';

const EnvironmentService = () => ({
  getRequest(environment: Environment): TRawEnvironment {
    return {
      type: 'Environment',
      spec: environment,
    };
  },

  validateDraft({name = '', description = '', values = []}: Environment) {
    return !!name && !!description && !!values.length;
  },
});

export default EnvironmentService();
