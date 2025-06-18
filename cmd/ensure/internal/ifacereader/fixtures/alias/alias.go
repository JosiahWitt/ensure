package alias

type ThisIsAnAlias = AnAliasPointsToThis

type AnAliasPointsToThis interface {
	Hello()
}
