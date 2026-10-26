export default {
  // Middleware
  EmailOrPasswordMissing: "E-mail a heslo jsou povinné",
  IncorrectCredentials: "Nesprávný e-mail nebo heslo, zkus to znovu",

  // Defaults
  InvitationCode: "",
  UserIdentifier: "user",
  CurrentQuizIdentifier: "user_played_quiz",
  CountTill: 3000,
  MinimumPoints: 0,
  MaximumPoints: 20,
  MinimumDurationInSeconds: 1,
  MaximumDurationInSeconds: 300,
  DefaultPoints: 1,
  DefaultDurationInSeconds: 10,
  ReadyMessage: "Připrav se…",
  AnswerSubmitted: "Odpověď odeslána",

  // Status
  Fail: "fail",
  Success: "success",
  Error: "error",

  // Quiz State
  Running: "running",
  Pause: "pause",

  // Components
  Question: "Question",
  Score: "Score",
  Waiting: "Waiting",

  // Events
  StartQuiz: "start_quiz",
  Authentication: "authentication",
  SentInvitaionCode: "send_invitation_code",
  GetQuestion: "send_question",
  Counter: "5_sec_counter",
  TerminateQuiz: "terminate_quiz",
  InvitationCodeValidation: "invitation_code_validation",
  RedirectToAdmin: "redirect_to_admin",
  NextQuestionAsk: "next_question",
  AdminDisconnected: "admin_is_disconnected",
  AskSkip: "ask_skip",
  AskForceSkip: "ask_force_skip",
  ShowScore: "show_score",
  AdminDisconnectedMessage: "Moderátor se odpojil, chvíli počkej…",
  SkipTimer: "skip_timer",
  EventAnswerSubmittedByUser: "answer submitted by user",
  EventPing: "ping",
  PauseQuiz: "pause_quiz",
  PauseQuizMessage: "Moderátor kvíz pozastavil, chvíli počkej",
  ResumeQuiz: "resume_quiz",
  ResumeQuizMessage: "Kvíz pokračuje, můžeš hrát dál",
  HostEndedQuizMessage: "Moderátor ukončil kvíz.",

  // Actions
  ActionAnserSubmittedByUser: "answer submitted by user",

  // Errors
  Unauthorized: "unauthorized to access resource",
  CodeNotFound: "invitation code not found",
  ReloadRequired: "Něco se pokazilo, obnov prosím stránku.",
  InvitationCodeNotFound: "invitation code not found",
  QuizSessionValidationFailed: "quiz-session-validation-failed",
  SessionWasCompleted: "session was completed",
  NoAnswerFound: "Vyber odpověď",

  // Messages
  CsvUploadSuccess: "Soubor byl nahrán",

  // Image upload
  validImageTypes: [
    "image/jpeg",
    "image/png",
    "image/gif",
    "image/webp",
    "image/heic",
    "image/heif",
  ],
};
