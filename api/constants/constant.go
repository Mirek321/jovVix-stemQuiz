package constants

// variables
const (
	CookieUser   = "user"
	KratosCookie = "ory_kratos_session"
)

// fiber contexts
const (
	ContextUid            = "userId"
	ContextUser           = "userContext"
	ContextQuizPermission = "quiz_permission"
)

// kratos
const (
	KratosID          = "kratosId"
	KratosUserDetails = "kratosUserDetails"
)

// params
const (
	ParamUid         = "userId"
	UserPlayedQuizId = "user_played_quiz_id"
	Username         = "username"
	MediaQuery       = "media"
	ParamTitle       = "title"
)

// Permissions
const (
	ReadPermission  = "read"
	WritePermission = "write"
	SharePermission = "share"
)

// Email templetes
const (
	QuizEmailSubject = "Máš pozvánku do kvízu"
)

// Success messages
// ...

// Fail messages
// ...
const (
	UsernameRequired = "Zadej uživatelské jméno"
)

// Error messages
const (
	ErrGetUser                     = "Nepodařilo se načíst uživatele"
	ErrLoginUser                   = "Přihlášení se nezdařilo"
	ErrInsertUser                  = "Nepodařilo se vytvořit uživatele, zkus to prosím později"
	ErrDeleteUser                  = "Nepodařilo se smazat uživatele, zkus to prosím později"
	ErrConvertTypeUser             = "Unable to convert to user-model type from locals"
	ErrHealthCheckDb               = "error while checking health of database"
	ErrUnauthenticated             = "Nepodařilo se ověřit tvou identitu"
	ErrUnauthorized                = "Přístup odepřen. Nemáš potřebná oprávnění."
	ErrShareNotAllowed             = "Nemáš oprávnění sdílet tento kvíz"
	ErrCannotChangeOwnAccess       = "Nemůžeš měnit svůj vlastní přístup k tomuto kvízu"
	ErrSharedQuizNotFound          = "Oprávnění ke sdílenému kvízu nebylo nalezeno"
	ErrGetSharedQuiz               = "Nepodařilo se načíst oprávnění ke sdílenému kvízu"
	ErrKratosAuth                  = "Nepodařilo se načíst uživatele z přihlašovací služby"
	ErrKratosDataInsertion         = "Nepodařilo se uložit údaje uživatele z přihlašovací služby"
	ErrKratosIDEmpty               = "Přihlášení nebylo nalezeno, přihlas se prosím znovu"
	ErrKratosCookieTime            = "Nepodařilo se ověřit platnost přihlášení"
	ErrRegisterQuiz                = "Nepodařilo se vytvořit kvíz"
	ErrCreatingDemoQuiz            = "Nepodařilo se vytvořit ukázkový kvíz"
	ErrQuizNotPublic               = "Tento kvíz není veřejný"
	ErrHideOrphanedQuiz            = "Kvíz bez autora nejde skrýt, nikdo by se k němu už nedostal"
	ErrGetTotalJoinUser            = "Nepodařilo se zjistit počet připojených uživatelů"
	ErrShareQuiz                   = "Nepodařilo se sdílet kvíz"
	ErrListShareQuiz               = "Nepodařilo se načíst seznam sdílených kvízů"
	ErrFetchAuthorizedUsersError   = "Nepodařilo se načíst uživatele s přístupem k vybranému kvízu."
	ErrCheckQuizCreatorExists      = "Nepodařilo se ověřit autora kvízu"
	ErrGetQuizPermission           = "Nepodařilo se načíst tvoje oprávnění ke kvízu"
	ErrUpdateUserPermissionForQuiz = "Nepodařilo se upravit oprávnění uživatele ke kvízu"
	ErrDeleteUserPermissionForQuiz = "Nepodařilo se odebrat oprávnění uživatele ke kvízu"
	ErrGetStreakCount              = "Nepodařilo se načíst počet sérií"
	ErrListCategories              = "Nepodařilo se načíst kategorie kvízů"
	ErrCreateCategory              = "Nepodařilo se vytvořit kategorii"
	ErrUpdateCategory              = "Nepodařilo se upravit kategorii"
	ErrDeleteCategory              = "Nepodařilo se smazat kategorii"
	ErrCategoryNotFound            = "Kategorie nebyla nalezena"
	ErrCategoryAlreadyExists       = "Kategorie s tímto názvem už existuje"
	ErrInvalidCoverImage           = "Titulní obrázek musí být obrázek"
	ErrCoverImageTooLarge          = "Titulní obrázek je příliš velký"
)

// Bad Request Message
const (
	BadRequestSharedQuizIdNotFound = "Chybí parametr shared_quiz_id"
)

// default Events
const (
	EventUserRegistered = "event:userRegistered"
)

// Middleware
const (
	// socket
	MiddlewareError     = "middleware_error"
	ErrorTypeConversion = "Nepodařilo se zpracovat data"

	// http/https
	ErrNotAllowed              = "Nemáš přístup k tomuto zdroji"
	ErrUserRequiredToCheckRole = "Nejsi přihlášený"

	// csv
	FileName                    = "file_name"
	MaxRows                     = 500
	FileSize                    = 100000 // TODO: change file size, ~1mb
	MaximumPoints               = 20
	MinimumPoints               = 0
	MaximumDurationInSeconds    = 300
	MinimumDurationInSeconds    = 1
	DefaultDurationInSeconds    = 30
	DefaultCsvPoints            = 1
	SheetName                   = "demo"
	QuizTitle                   = "quiz_title"
	QuizTitleRequired           = "Zadej název kvízu."
	QuizTitleMaxLength          = 50
	ErrQuizTitleTooLong         = "Název kvízu je příliš dlouhý. Zkrať ho prosím (maximálně 50 znaků)."
	ErrGettingAttachment        = "Nepodařilo se načíst soubor"
	ErrFileSizeExceed           = "Soubor je příliš velký. Vyber prosím menší soubor."
	ErrFileIsNotInSupportedType = "Nepodporovaný typ souboru"
	ErrProblemInUploadFile      = "Při nahrávání souboru došlo k chybě"
	ErrValidatingColumns        = "Sloupce v souboru nemají správný formát"
	ErrParsingFile              = "Soubor se nepodařilo zpracovat"
	ErrRowsReachesToMaxCount    = "Soubor má příliš mnoho řádků"
	ErrSurveyAnswerLength       = "u ankety musí být vyplněna aspoň jedna správná odpověď"
	ErrSingleAnswerLength       = "u otázky s jednou odpovědí musí být právě jedna správná odpověď"
	ErrQuestionType             = "zadej platný typ otázky"
	ErrQuestionId               = "Tento typ otázky neexistuje"
	ErrAppendQuestionsFailed    = "Otázky se nepodařilo přidat do kvízu"
	ErrEmptyFile                = "Nahraný soubor je prázdný. Vyber soubor s obsahem."
	ErrUnsupportedFileType      = "Nahraný soubor není platný soubor CSV. Zkontroluj formát a zkus to znovu."
	ErrEmptyQuestionText        = "chybí text otázky"
	ErrInsufficientOptions      = "jsou potřeba aspoň 2 možnosti"
	ErrEmptyCorrectAnswer       = "chybí správná odpověď"
	ErrInvalidCorrectAnswer     = "správná odpověď musí být číslo existující možnosti"
	ErrInvalidPoints            = "body musí být celé číslo v povoleném rozsahu"
	ErrInvalidCSVRows           = "Soubor CSV obsahuje chybné řádky, oprav je a zkus to znovu."
	ErrInvalidQuestionTimeLimit = "Časový limit otázek není správně nastavený"
	ErrInvalidQuestionMedia     = "typ obsahu otázky musí být jeden z: text, image, code"
	ErrInvalidOptionsMedia      = "typ obsahu možností musí být jeden z: text, image, code"

	// quiz-id
	QuizId       = "quiz_id"
	QuestionId   = "question_id"
	SharedQuizId = "shared_quiz_id"
	CategoryId   = "category_id"

	// Base64 data-URI cover images inflate ~4/3 over the raw file; 1 MiB of
	// text comfortably covers the 500 KB client-side file limit.
	MaxCoverImageBytes = 1 << 20
)

var AllowedCoverImageTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/gif":  true,
	"image/webp": true,
	"image/heic": true,
	"image/heif": true,
}

// components
const (
	Waiting  = "Waiting"
	Question = "Question"
	Score    = "Score"
	Loading  = "Loading"
	Running  = "Running"

	ToUser  = 1
	ToAdmin = 2
	ToAll   = 3
)

// constants
const (
	MinInvitationCode = 100000
	MaxInvitationCode = 999999
	Counter           = 5
	Count             = 3
	StreakBaseScore   = 100
)

// Quiz Events
const (
	// Event 1. Authentication <admin side>
	EventAuthentication  = "authentication"
	ActionAuthentication = "authentication check"
	Unauthenticated      = "Pro přístup se musíš přihlásit"
	InvalidCredentials   = "Neplatné přihlašovací údaje"

	// Event 2. Authorization <admin/user side>
	EventAuthorization  = "authorization"
	ActionAuthorization = "check for access"
	UserNotExist        = "Uživatel neexistuje"
	Unauthorized        = "unauthorized to access resource"

	// Event 3. Session Validation <admin>
	EventSessionValidation  = "session_validation"
	ActionSessionValidation = "session validation from server side"
	ErrSessionNotFound      = "Sezení není dostupné"

	// Event 4. UserSession Validation <admin/user>
	EventUserSessionValidation   = "user_validation"
	ActionUserSessionValidation  = "user session get or create"
	CurrentUserQuiz              = "user_played_quiz"               // use by web
	ErrUserQuizSessionValidation = "quiz-session-validation-failed" // use by web
	ErrAdminCannotBeUser         = "Moderátor nemůže hrát ve vlastním kvízu"

	EventRedirectToAdmin     = "redirect_to_admin"
	ActionCurrentUserIsAdmin = "current user is admin"

	// Event 4. Active session <admin>
	EventActivateSession          = "session_activation"
	EventSendInvitationCode       = "send_invitation_code" // use by web
	ActionSessionActivation       = "activate demanded session and sent invitation code"
	QuizSessionInvitationCode     = "invitationCode"
	SessionIDParam                = "session_id"
	ActiveQuizObj                 = "current active quiz obj"
	NoPlayerFound                 = "Zatím se nepřipojil žádný hráč"
	StartQuizByAdminNoPlayerFound = "start quiz by admin but no player found"
	ActionSendUserData            = "send user join data"
	JoinUserOnRunningQuiz         = "join_user_on_running_quiz"

	// Event 5. Join quiz <User>
	EventJoinQuiz                  = "invitation_code_validation"
	ActionJoinQuiz                 = "invitation code validation"
	ErrInvitationCodeInWrongFormat = "Kód kvízu nemá správný formát"
	ErrInvitationCodeNotFound      = "invitation code not found" // use by web
	ErrSessionWasCompleted         = "session was completed"     // use by web
	ErrMaxTryToGenerateCode        = "Nepodařilo se vygenerovat kód kvízu, zkus to znovu"

	UserName          = "username"
	UserUkey          = "users_username_ukey"
	Join              = "join access"
	QuizStartsSoon    = "Kvíz brzy začne"
	ErrUsernameExists = "Toto uživatelské jméno už existuje"

	// Event 6. Start quiz <admin>
	EventUserJoined       = "user joined"
	EventStartQuiz        = "start_quiz"       // use by web
	EventStartQuizByAdmin = "startQuizByAdmin" // use by admin for start quiz

	// Event 7. Get Questions
	EventSendQuestion              = "send_question"
	ActionSendQuestion             = "send single question to user"
	QuizQuestionStatus             = "quiz question status"
	GetQuestions                   = "get quiz questions"
	NextQuestionWillServeSoon      = "Další otázka už brzy"
	ErrInGettingQuestion           = "error during getting question"
	ErrInGettingTotalQuestionCount = "error during getting total question count"

	EventPublishQuestion       = "publish_question"
	EventStartCount5           = "5_sec_counter" // use by web
	ActionCounter              = "5 second counter"
	EventNextQuestionAsked     = "next_question"         // use by web
	AdminDisconnected          = "admin_is_disconnected" // use by web
	EventAnswerSubmittedByUser = "answer submitted by user"
	ActionAnserSubmittedByUser = "answer submitted by user"

	// Event 8. Submit answer
	ErrQuizNotFound           = "Kvíz nebyl nalezen"
	ErrAnswerSubmit           = "Neplatné údaje odpovědi"
	ErrAnswerAlreadySubmitted = "Odpověď už byla odeslána"
	ErrQuestionNotActive      = "Na tuto otázku už nejde odpovídat"
	ErrPublishAnswer          = "error while publishing answer in redis"

	// Event skip
	EventSkipAsked  = "ask_skip" // use by web
	WarnSkip        = "Někteří hráči ještě neodpověděli. Chceš otázku přeskočit?"
	EventForceSkip  = "ask_force_skip"
	EventSkipTimer  = "skip_timer"
	EventPauseQuiz  = "pause_quiz"
	EventResumeQuiz = "resume_quiz"

	// Event 8. Get score page
	EventShowScore  = "show_score"
	ActionShowScore = "show score page during quiz"

	// Event 9. Terminate quiz
	EventTerminateQuiz  = "terminate_quiz"
	ActionTerminateQuiz = "terminate quiz after completing"

	ErrActiveDeleteQuiz = "Tuto akci nelze provést, dokud kvíz probíhá."

	// Event 10. unhandled event
	UnknownError  = "unknown_error"
	ErrJWTExpired = "Platnost tokenu vypršela, zkus to prosím znovu"

	// Event 11. ping
	EventPing = "ping"
	EventPong = "pong"
)

// final scoreboard cookie for user
const UserPlayedQuiz = "user_played_quiz"
const ActiveQuizId = "active_quiz_id"

// database table names
const (
	UserQuizResponsesTable   = "user_quiz_responses"
	UserPlayedQuizzesTable   = "user_played_quizzes"
	QuestionsTable           = "questions"
	UsersTable               = "users"
	ActiveQuizzesTable       = "active_quizzes"
	QuizQuestionsTable       = "quiz_questions"
	ActiveQuizQuestionsTable = "active_quiz_questions"
	QuizzesTable             = "quizzes"
)

// Question Types
const (
	SingleAnswerString = "single answer"
	SurveyString       = "survey"

	SingleAnswer = 1
	Survey       = 2
)

// Media Types
const (
	MediaText  = "text"
	MediaImage = "image"
	MediaCode  = "code"
)

// Pagination and Filters

const (
	PageNumberQueryParam = "page"
	NameQueryParam       = "name"
	OrderQueryParam      = "order"
	OrderByQueryParam    = "orderBy"
	DefaultPageSize      = 10
)

// Channel name for redis pubsub
const (
	ChannelUserJoin       = "user_joined"
	ChannelUserDisconnect = "user_disconnect"
	ChannelSetAnswer      = "set_answer"
)

// Session modes
const (
	ModeLive      = "live"
	ModeSelfPaced = "self_paced"

	// SelfPacedGracePeriodSeconds tolerates network latency for late self-paced answers.
	SelfPacedGracePeriodSeconds = 2

	ModeQueryParam = "mode"

	EventSelfPacedSession = "self_paced_session"

	ErrInvalidSessionMode          = "Neplatný režim spuštění kvízu."
	ErrSessionNotSelfPaced         = "Tohle sezení se nehraje vlastním tempem."
	ErrSelfPacedQuizEnded          = "Kvíz už skončil."
	ErrSelfPacedQuestionNotActive  = "Tahle otázka teď není aktivní."
	ErrSelfPacedSessionAdminSocket = "Tenhle kvíz se hraje vlastním tempem."
	ErrNotParticipant              = "Odpovídat můžou jen účastníci sezení"
	ErrSessionStatusForbidden      = "Ke stavu tohohle sezení nemáš přístup."
	ErrQuizHasNoQuestions          = "Kvíz nemá žádné otázky. Nejdřív nějakou přidej."
)
