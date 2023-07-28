export interface OTELYaml {
  groups: Group[];
}

export interface Group {
  id: string;
  prefix: string;
  type: string;
  brief: string;
  note?: string;
  attributes: Attribute[];
  constraints?: Constraint[];
  extends?: string;
  span_kind?: string;
}

export interface Attribute {
  id?: string;
  type?: TypeClass | TypeEnum;
  requirement_level?: RequirementLevelClass | string;
  brief?: string;
  sampling_relevant?: boolean;
  examples?: Array<number | string> | number | string;
  note?: string;
  ref?: string;
}
export interface CompleteAttribute extends Attribute {
  group: string;
}

export interface RequirementLevelClass {
  conditionally_required?: string;
  recommended?: string;
}

export interface TypeClass {
  allow_custom_values: boolean;
  members: Member[];
}

export interface Member {
  id: string;
  value: string;
  brief: string;
}

export enum TypeEnum {
  Int = 'int',
  String = 'string',
}

export interface Constraint {
  include: string;
}
