package router

import (
	v1 "github.com/auroraride/aurservd/app/router/v1"
	v2 "github.com/auroraride/aurservd/app/router/v2"
)

func loadManagerRoutes() {
	v1.LoadManagerV1Routes(root)
	v2.LoadManagerV2Routes(root)
}
