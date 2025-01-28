package text

var En = [...]string{
//Default
	"%s",

//Join
	"Join",
//Create
	"Create",

//GameStarted
	"Game started!",
//GameNotSelected
	"",
//GameMafia
	"Mafia",

//IsMafiaF
	"%s is mafia!",
//IsNotMafiaF
	"%s is not mafia!",
//HealedF
	"You healed %s!",
//ChosenToKillF
	"You chose to kill %s!",
//NightFalls
	"Night falls. Civilians fall asleep",
//DayIsComing
	"The day is coming. Civilians are waking up",
//IsWakingUpF
	"%s is waking up!",
//FallAsleepF
	"%s made choice and fell asleep!",
//MakeChoiceF
	"[%s] %s, make your choice:",
//MafiaSuccessF
	"Bad news!\nMafia killed %s!",
//MafiaFailed
	"Good news!\nThis night doctor saved victim's life and no one died!",

//MafiaWon
	"Mafia won!",
//CiviliansWon
	"Civilians won!",

//Voting
	"Vote who you want to kick:",
//VotedF
	"You voted for %s",
//VotingResultsF
	"Voting results:\n%s",
//VotesAreEqual
	"Votes are equal! Vote again",
//VoteKickF
	"%s is kicked! He was %s",

//GameChosenF
	"You chose %s.\nTo start game send /game_start",
//CantStartGame	
	"You can't start game until you choose game type!",

//AtLeastMembersF
	"You need to have at least %d members in lobby to start game!",
//RoleF
	"You are %s!",

//KeyCreatedF
	"Your lobby key is: %s\nNow choose what game you want to play!",
//SendKey
	"Write key for lobby to join, like: O1F5",
//IncorrectKey
	"Incorrect key!\nMake sure you use the correct key",
//LobbyNotExists
	"Lobby with that key does not exist or your key is incorrect!",
//LobbyGameIsStarted
	"In that lobby the game has already started! You can't join",

//PlayerJoinedLobbyF
	"%s joined the lobby!\nCurrent members in lobby:\n%s",
//AlreadyInLobbyF
	"You are currently in a lobby with key: %s!\nLeave that lobby and then you can join",

//UnknownCommand
	"I don't know that command",
//StartCommandF
	"Hello %s, I'm Games for fun bot!\nChoose to join a lobby or create one to start the game!",

//CreatingLobbyError
	"Something went wrong while creating the lobby! Please try again later.",
//CantFindUser
	"Error!\nCan't find you in data! Please restart GFFBot",
//SomethingWentWrong
	"Something went wrong! Please restart the bot",

//You
	"You",

//Yes
	"Yes",
//No
	"No",

//SugesToKillF
	"%s sugest to choose %s as victim",
//Accepted
	"%s accepted",
//Declined
	"%s declined",
}