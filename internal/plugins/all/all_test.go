package all_test

import (
	"testing"

	"github.com/JosiahWitt/ensure"
	"github.com/JosiahWitt/ensure/internal/plugins/all"
	"github.com/JosiahWitt/ensure/internal/plugins/mocks"
	"github.com/JosiahWitt/ensure/internal/plugins/setupmocks"
	"github.com/JosiahWitt/ensure/internal/plugins/subject"
)

func TestTablePlugins(t *testing.T) {
	ensure := ensure.New(t)

	tablePlugins := all.TablePlugins()
	ensure(len(tablePlugins)).Equals(3)

	_, ok0 := tablePlugins[0].(*mocks.TablePlugin)
	ensure(ok0).IsTrue()

	_, ok1 := tablePlugins[1].(*setupmocks.TablePlugin)
	ensure(ok1).IsTrue()

	_, ok2 := tablePlugins[2].(*subject.TablePlugin)
	ensure(ok2).IsTrue()
}
