package anonymizer

import (
	"context"
	"io"
	"strconv"
	"time"

	"github.com/docker/go-units"
	"github.com/nixys/nxs-data-anonymizer/ctx"
	"github.com/nixys/nxs-data-anonymizer/misc"
	"github.com/sirupsen/logrus"

	appctx "github.com/nixys/nxs-go-appctx/v3"
)

type anonymizeOpts struct {
	c  context.Context
	l  *logrus.Logger
	ch chan error
	db ctx.DBCtx
	w  io.Writer
	a  misc.Anonymizer
}

func Runtime(app appctx.App) error {

	var (
		// Bytes count printed in log last time
		lb int64

		// Timer to print log read bytes count
		timer *time.Timer
	)

	cc := app.ValueGet().(*ctx.Ctx)

	cx, cf := context.WithCancel(app.SelfCtx())
	defer cf()

	ch := make(chan error, 1)

	timer = time.NewTimer(cc.Progress.Rhythm)
	if cc.Progress.Rhythm == 0 {
		timer.Stop()
	}

	if err := anonymize(
		anonymizeOpts{
			c:  cx,
			l:  cc.Log,
			ch: ch,
			db: cc.DB,
			w:  cc.Output,
			a:  cc.Anonymizer,
		},
	); err != nil {
		return err
	}

	for {
		select {
		case <-app.SelfCtxDone():

			// Log reader progress if necessary
			if cc.Progress.Rhythm != 0 && lb != cc.PR.Bytes() {
				progressLog(cc.Log, cc.PR.Bytes(), cc.Progress.Humanize)
			}

			cc.Log.Info("anonymizer routine done")
			return nil
		case err := <-ch:

			// Log reader progress if necessary
			if cc.Progress.Rhythm != 0 && lb != cc.PR.Bytes() {
				progressLog(cc.Log, cc.PR.Bytes(), cc.Progress.Humanize)
			}

			if err != nil {

				cc.Log.WithFields(logrus.Fields{
					"details": err,
				}).Errorf("anonymize")

				return err
			}

			cc.Log.Info("anonymizer routine done")
			return nil
		case <-timer.C:

			// Save bytes count printed in log last time
			lb = cc.PR.Bytes()

			// Log reader progress
			progressLog(cc.Log, lb, cc.Progress.Humanize)

			timer.Reset(cc.Progress.Rhythm)
		}
	}
}

func anonymize(st anonymizeOpts) error {

	if st.db.Type == ctx.DBTypeMySQL && st.db.Cleanup == true && st.db.MySQL != nil {
		if err := st.db.MySQL.DBCleanup(); err != nil {

			st.l.WithFields(logrus.Fields{
				"details": err,
			}).Errorf("anonymize: MySQL clean up")

			return err
		}
	}

	go func() {
		st.ch <- st.a.Run(st.c, st.w)
	}()

	return nil
}

func progressLog(l *logrus.Logger, b int64, h bool) {

	var s string

	// Prepare output bytes string
	if h == true {
		s = units.BytesSize(float64(b))
	} else {
		s = strconv.FormatInt(b, 10)
	}

	l.WithFields(
		logrus.Fields{
			"read bytes": s,
		},
	).Info("anonymization progress")
}
