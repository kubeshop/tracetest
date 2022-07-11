import {VariableDefinition} from 'postman-collection';

export function substituteVariable(variables: VariableDefinition[], value: string | undefined) {
  const regExpMatchArray = (value || '').match(/\{{([^}]+)}}/);
  return regExpMatchArray
    ? value
        ?.replaceAll(
          regExpMatchArray?.[0],
          variables.find(variable => variable.key === regExpMatchArray?.[1])?.value || value
        )
        .replaceAll(' ', '')
    : value;
}
