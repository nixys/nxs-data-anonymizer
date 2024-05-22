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

	mysql_anonymize "github.com/nixys/nxs-data-anonymizer/modules/anonymizers/mysql"
	pgsql_anonymize "github.com/nixys/nxs-data-anonymizer/modules/anonymizers/pgsql"
	"github.com/nixys/nxs-data-anonymizer/modules/filters/relfilter"
	progressreader "github.com/nixys/nxs-data-anonymizer/modules/progress_reader"
)

type anonymizeSettings struct {
	c  context.Context
	l  *logrus.Logger
	pr *progressreader.ProgressReader
	ch chan error
	db ctx.DBCtx
	rs relfilter.Rules
	w  io.Writer
	s  ctx.SecurityCtx
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

	// Init progress reader
	pr := progressreader.Init(cc.Input)

	c := make(chan error, 1)

	timer = time.NewTimer(cc.Progress.Rhythm)
	if cc.Progress.Rhythm == 0 {
		timer.Stop()
	}

	if err := anonymize(
		anonymizeSettings{
			c:  cx,
			l:  cc.Log,
			pr: pr,
			ch: c,
			db: cc.DB,
			rs: cc.Rules,
			w:  cc.Output,
			s:  cc.Security,
		},
	); err != nil {
		return err
	}

	for {
		select {
		case <-app.SelfCtxDone():

			// Log reader progress if necessary
			if cc.Progress.Rhythm != 0 && lb != pr.Bytes() {
				progressLog(cc.Log, pr.Bytes(), cc.Progress.Humanize)
			}

			cc.Log.Info("anonymizer routine done")
			return nil
		case err := <-c:

			// Log reader progress if necessary
			if cc.Progress.Rhythm != 0 && lb != pr.Bytes() {
				progressLog(cc.Log, pr.Bytes(), cc.Progress.Humanize)
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
			lb = pr.Bytes()

			// Log reader progress
			progressLog(cc.Log, lb, cc.Progress.Humanize)

			timer.Reset(cc.Progress.Rhythm)
		}
	}
}

func anonymize(st anonymizeSettings) error {

	// Anonymizer reader
	var ar io.Reader

	// Init anonymize reader in accordance with specified database type
	switch st.db.Type {
	case ctx.DBTypeMySQL:

		// Drop database tables if necessary (experimental)
		if st.db.Cleanup == true && st.db.MySQL != nil {
			if err := st.db.MySQL.DBCleanup(); err != nil {

				st.l.WithFields(logrus.Fields{
					"details": err,
				}).Errorf("anonymize: MySQL clean up")

				return err
			}
		}

		ar = mysql_anonymize.Init(
			st.c,
			st.pr,
			mysql_anonymize.InitSettings{
				Security: mysql_anonymize.SecuritySettings{
					TablePolicy:     st.s.TablePolicy,
					TableExceptions: st.s.TableExceptions,
				},
				Rules: st.rs,
			},
		)
	case ctx.DBTypePgSQL:
		ar = pgsql_anonymize.Init(st.c, st.pr, st.rs)
	default:

		st.l.WithFields(logrus.Fields{
			"details": "unknown database type",
		}).Errorf("anonymize")

		return misc.ErrRuntime
	}

	go func() {
		_, err := io.Copy(st.w, ar)
		st.ch <- err
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
