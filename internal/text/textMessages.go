package text

const (
	LettersBytes = "QWERTYUIOPASDFGHJKLZXCVBNM"
	DigitsBytes  = "1234567890"
)

const (
	Default = iota

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
)
