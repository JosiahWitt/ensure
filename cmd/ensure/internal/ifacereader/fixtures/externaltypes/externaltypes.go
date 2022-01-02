package externaltypes

import (
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example1"
	"github.com/JosiahWitt/ensure/cmd/ensure/internal/ifacereader/fixtures/externaltypes/example2"
)

type ExternalTypes interface {
	ExternalIO(a map[example2.Float64]*example1.Message) map[example1.String]*example2.User
}
