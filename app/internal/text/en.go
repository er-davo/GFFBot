package text

var En = [...]string{
	Civilian:  "Civilian",
	Mafia:     "Mafia",
	Detective: "Detective",
	Doctor:    "Doctor",

	GMafia:  "Mafia",
	GBunker: "Shelter",

	Default: "%s",

	Join:   "Join",
	Create: "Create",

	GameStarted:     "Game started!",
	GameNotSelected: "",
	GameMafia:       "Mafia",

	IsMafiaF:      "%s is mafia!",
	IsNotMafiaF:   "%s is not mafia!",
	HealedF:       "You healed %s!",
	ChosenToKillF: "You chose to kill %s!",
	NightFalls:    "Night falls. Civilians fall asleep",
	DayIsComing:   "The day is coming. Civilians are waking up",
	IsWakingUpF:   "%s is waking up!",
	FallAsleepF:   "%s made choice and fell asleep!",
	MakeChoiceF:   "[%s] %s, make your choice:",
	MafiaSuccessF: "Bad news!\nMafia killed %s!",
	MafiaFailed:   "Good news!\nThis night doctor saved victim's life and no one died!",

	MafiaWon:     "Mafia won!",
	CiviliansWon: "Civilians won!",

	Voting:         "Vote who you want to kick:",
	VotedF:         "You voted for %s",
	VotingResultsF: "Voting results:\n%s",
	VotesAreEqual:  "Votes are equal! Vote again",
	VoteKickF:      "%s is kicked! He was %s",

	GameChosenF:   "You chose %s.\nTo start game send /game_start",
	CantStartGame: "You can't start game until you choose game type!",

	AtLeastMembersF: "You need to have at least %d members in lobby to start that game!",
	MaximumMembersF: "You can have maximum %d members in lobby to start that game!",
	RoleF:           "You are %s!",

	KeyCreatedF:        "Your lobby key is: %s\nNow choose what game you want to play!",
	SendKey:            "Write key for lobby to join, like: O1F5",
	IncorrectKey:       "Incorrect key!\nMake sure you use the correct key",
	LobbyNotExists:     "Lobby with that key does not exist or your key is incorrect!",
	LobbyGameIsStarted: "In that lobby the game has already started! You can't join",

	PlayerJoinedLobbyF: "%s joined the lobby!\nCurrent members in lobby:\n%s",
	AlreadyInLobbyF:    "You are currently in a lobby with key: %s!\nLeave that lobby and then you can join",

	UnknownCommand:    "I don't know that command",
	StartCommandF:     "Hello %s, I'm the Games for Fun bot!\nYou can play games like Mafia or Shelter.\nUse /help for more information.",
	HelpCommand:       "Available commands:\n/login - log in to save your game statistics.\n/statistic - view your game statistics\n/lobby - create or join a lobby.\n/game_start - start the game",
	LoginCommand:      "Logging in or signing up...",
	StatisticCommandF: "Your statistics:\n%s",
	LobbyCommand:      "Create a lobby or join an existing one:",
	
	LoginSuccess: "You have successfully logged in!",
	LoginAlready: "You are already logged in!",
	LoginFailed:  "Login failed! Please try again later.",

	CreatingLobbyError: "Something went wrong while creating the lobby! Please try again later.",
	CantFindUser:       "Error!\nCan't find you in data! Please restart GFFBot with /start",
	SomethingWentWrong: "Something went wrong! Please restart the bot",

	You: "You",

	Yes: "Yes",
	No:  "No",

	SugestToKillF: "%s sugest to choose %s as victim",
	AcceptedF:     "%s accepted",
	DeclinedF:     "%s declined",

	GameStoped: "Something went wrong. Game stoped, start new game.",
}
