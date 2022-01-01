package externaltypes

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example1"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example2"
)

type ExternalTypes interface {
	ExternalIO(a *example1.Message) *example2.User
}
