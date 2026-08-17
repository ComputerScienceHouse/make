type BaseQuestion = {
  id: number
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

export type QuestionWithAnswer =
  | (TextQuestion & { answer: string })
  | (BigTextQuestion & { answer: string })
  | (NumberQuestion & { answer: number })
  | (RadioQuestion & { answer: string })
  | (CheckboxQuestion & { answer: boolean })

export interface Training {
  id: number
  title: string
  description: string
  requiredCorrect: number
  showAnswers: boolean
  questions: Question[]
}

export interface TrainingFull {
  id: number
  title: string
  description: string
  requiredCorrect: number
  showAnswers: boolean
  questions: QuestionWithAnswer[]
}

export interface UserTraining {
  userUuid: string
  trainingId: number
  completedAt: Date
  expiresAt: Date
}

export interface SubmissionResults {
  passed: boolean
  numCorrect: number
  numIncorrect: number
  grade: number
  gradedResponse: Record<string, boolean>
}
