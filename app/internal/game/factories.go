package game

type GameFactory interface {
	CreateGame(isStarted *bool, members *Users) GameStarter
	CreatePlayer(lang *string) Player
}

type MafiaGameFactory struct{}

func (m MafiaGameFactory) CreateGame(isStarted *bool, members *Users) GameStarter {
    return &MafiaGame{
		IsStarted: isStarted,
		Members: members,
	}
}

func (m MafiaGameFactory) CreatePlayer(lang *string) Player {
    return &MafiaPlayer{Lang: lang}
}

type BunkerGameFactory struct{}

func (b BunkerGameFactory) CreateGame(isStarted *bool, members *Users) GameStarter {
	return &BunkerGame{
        IsStarted: isStarted,
        Members: members,
    }
}

func (b BunkerGameFactory) CreatePlayer(lang *string) Player {
    return &BunkerPlayer{Lang: lang}
}