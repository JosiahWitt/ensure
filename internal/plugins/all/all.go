// Package all combines all the plugins in the correct order.
package all

import (
	"github.com/JosiahWitt/ensure/internal/plugins"
	"github.com/JosiahWitt/ensure/internal/plugins/internal/mocks"
	mocksplugin "github.com/JosiahWitt/ensure/internal/plugins/mocks"
	"github.com/JosiahWitt/ensure/internal/plugins/setupmocks"
	"github.com/JosiahWitt/ensure/internal/plugins/subject"
)

// TablePlugins provides all the plugins for table-driven tests in the correct order.
func TablePlugins() []plugins.TablePlugin {
	m := &mocks.All{}

	// This order matters. Mocks are loaded and setup, and then the subject is populated.
	return []plugins.TablePlugin{
		mocksplugin.New(m),
		setupmocks.New(),
		subject.New(m),
	}
}
