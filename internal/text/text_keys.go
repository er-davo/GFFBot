package text

const (
	LettersBytes = "QWERTYUIOPASDFGHJKLZXCVBNM"
	DigitsBytes  = "1234567890"
)

const (
	Civilian = iota
	Mafia
	Detective
	Doctor

	GMafia
	GBunker

	Profession
	BoilogicalParams
	HealthStatus
	Hobby
	Phopia
	Character
	Skill
	Knowledge
	Baggage
	ActionCard
	ConditionCard

	Default

	Join
	Create

	GameStarted
	GameNotSelected
	GameMafia

	IsMafiaF
	IsNotMafiaF
	HealedF
	ChosenToKillF

	NightFalls
	DayIsComing
	IsWakingUpF
	FallAsleepF
	MakeChoiceF

	MafiaSuccessF
	MafiaFailed

	MafiaWon
	CiviliansWon

	Voting
	VotedF
	VotingResultsF
	VotesAreEqual
	VoteKickF

	GameChosenF
	CantStartGame

	AtLeastMembersF
	RoleF

	KeyCreatedF
	SendKey
	IncorrectKey
	LobbyNotExists
	LobbyGameIsStarted

	PlayerJoinedLobbyF
	AlreadyInLobbyF

	UnknownCommand
	StartCommandF

	CreatingLobbyError
	CantFindUser
	SomethingWentWrong

	You

	Yes
	No

	SugestToKillF
	AcceptedF
	DeclinedF

	GameStoped
)
