type BaseQuestion = {
  id: string
  label: string
  required?: boolean
}

export type TextQuestion = BaseQuestion & {
  type: 'text'
}

export type BigTextQuestion = BaseQuestion & {
  type: 'textarea'
}

export type NumberQuestion = BaseQuestion & {
  type: 'number'
}

export type RadioQuestion = BaseQuestion & {
  type: 'radio'
  options: string[]
}

export type CheckboxQuestion = BaseQuestion & {
  type: 'checkbox'
}

export type Question =
  | TextQuestion
  | BigTextQuestion
  | NumberQuestion
  | RadioQuestion
  | CheckboxQuestion

export interface Training {
  id: number
  title: string
  description: string
  questions: Question[]
}
export interface UserTraining {
  userUuid: string
  trainingId: number
  completedAt: Date
  expiresAt: Date
}
