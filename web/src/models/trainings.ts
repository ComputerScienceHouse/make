export interface Training {
    id: number,
    title: string,
    description: string,
    questions: string,
}
export interface UserTraining {
  userUuid: string;
  trainingId: number;
  completedAt: Date;
  expiresAt: Date;
}