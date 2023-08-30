package anonymizer

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/docker/go-units"
	"github.com/nixys/nxs-data-anonymizer/ctx"
	"github.com/sirupsen/logrus"

	appctx "github.com/nixys/nxs-go-appctx/v2"

	mysql_anonymize "github.com/nixys/nxs-data-anonymizer/modules/anonymizers/mysql"
	pgsql_anonymize "github.com/nixys/nxs-data-anonymizer/modules/anonymizers/pgsql"
	progressreader "github.com/nixys/nxs-data-anonymizer/modules/progress_reader"
)

// Runtime executes the routine
func Runtime(cr context.Context, appCtx *appctx.AppContext, crc chan interface{}) {

	var (
		// Anonymizer reader
		ar io.Reader

		// Bytes count printed in log last time
		lb int64

		// Timer to print log read bytes count
		timer *time.Timer
	)

	cc := appCtx.CustomCtx().(*ctx.Ctx)

	crr, cf := context.WithCancel(cr)

	// Init progress reader
	pr := progressreader.Init(cc.Args.Input)

	// Init anonymize reader in accordance with specified database type
	switch cc.Args.DBType {
	case ctx.DBTypeMySQL:

		// Drop database tables if necessary (experimental)
		if cc.Args.Cleanup == true && cc.MySQL != nil {
			if err := cc.MySQL.DBCleanup(); err != nil {
				appCtx.Log().Errorf("MySQL clean up error: %s", err)
				appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
				cf()
				return
			}
		}

		ar = mysql_anonymize.Init(crr, pr, cc.Rules)
	case ctx.DBTypePgSQL:
		ar = pgsql_anonymize.Init(crr, pr, cc.Rules)
	default:
		appCtx.Log().Error("unknown database type")
		appCtx.RoutineDoneSend(appctx.ExitStatusFailure)
		cf()
		return
	}

	c := make(chan error, 1)

	timer = time.NewTimer(cc.Progress.Rhythm)
	if cc.Progress.Rhythm == 0 {
		timer.Stop()
	}

	go func() {
		_, err := io.Copy(cc.Args.Output, ar)
		c <- err
	}()

	for {
		select {
		case <-cr.Done():
			// Program termination.

			// Log reader progress if necessary
			if cc.Progress.Rhythm != 0 && lb != pr.Bytes() {
				progressLog(appCtx.Log(), pr.Bytes(), cc.Progress.Humanize)
			}

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
		case <-timer.C:

			// Save bytes count printed in log last time
			lb = pr.Bytes()

			// Log reader progress
			progressLog(appCtx.Log(), lb, cc.Progress.Humanize)

			timer.Reset(cc.Progress.Rhythm)
		}
	}
}

func progressLog(l *logrus.Logger, b int64, h bool) {

	var s string

	// Prepare output bytes string
	if h == true {
		s = units.HumanSize(float64(b))
	} else {
		s = strconv.FormatInt(b, 10)
	}

	l.WithFields(
		logrus.Fields{
			"read bytes": s,
		},
	).Info("anonymization progress")
}
