export enum ResourceType {
  Test = 'tests',
  TestSuite = 'testsuites',
  VariableSet = 'variablesets',
}

export const ResourceName = {
  [ResourceType.Test]: 'Test',
  [ResourceType.TestSuite]: 'Test Suite',
  [ResourceType.VariableSet]: 'Variable Set',
} as const;
