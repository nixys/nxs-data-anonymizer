package anonymizer

import (
	"context"
	"example/ctx"
	"io"

	appctx "github.com/nixys/nxs-go-appctx/v2"

	mysql_anonymize "example/modules/anonymizers/mysql"
	pgsql_anonymize "example/modules/anonymizers/pgsql"
)

// Runtime executes the routine
func Runtime(cr context.Context, appCtx *appctx.AppContext, crc chan interface{}) {

	var r io.Reader

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	crr, cf := context.WithCancel(cr)

	switch cc.Args.DBType {
	case ctx.DBTypeMySQL:

		// Drop database tables if necessary
		if cc.Args.Cleanup == true && cc.MySQL != nil {
			if err := cc.MySQL.DBCleanup(); err != nil {
				appCtx.Log().Errorf("MySQL clean up error: %s", err)
				appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
				cf()
				return
			}
		}

		r = mysql_anonymize.Init(crr, cc.Args.Input, cc.Rules)
	case ctx.DBTypePgSQL:
		//appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
		r = pgsql_anonymize.Init(crr, cc.Args.Input, cc.Rules)
	default:
		appCtx.Log().Error("unknown database type")
		appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
		cf()
		return
	}

	c := make(chan error, 1)

	go func() {
		_, err := io.Copy(cc.Args.Output, r)
		c <- err
	}()

	for {
		select {
		case <-cr.Done():
			// Program termination.
			appCtx.Log().Info("anonymizer routine done")
			cf()
			return
		case <-crc:
			// Updated context application data.
			// Set the new one in current goroutine.
			appCtx.Log().Info("anonymizer routine reload, ignoring")
		case err := <-c:
			if err != nil {
				appCtx.Log().Errorf("anonymizer routing error: %s", err)
				appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
			} else {
				appCtx.RoutineDoneSend(appctx.ExitStatusSuccess)
			}
		}
	}
}
