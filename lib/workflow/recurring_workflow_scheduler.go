package workflow

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"
	rrule "github.com/teambition/rrule-go"

	"github.com/sirupsen/logrus"
)

// maxWakeupInterval is the maximum amount of time we will sleep
// between checking the recurrance table.
var maxWakeupInterval = 1 * time.Minute

type Schedule struct {
	// NOTE(ssd) 2019-05-13: Since name and workflow-name are
	// user-controlled in the case of many scheduled workflows, we
	// need the ID to create unique workflow names.
	ID           int64
	Enabled      bool
	Name         string
	WorkflowName string
	Parameters   []byte
	Recurrence   string
	NextDueAt    time.Time
}

type workflowScheduler struct {
	backend Backend
}

func (w *workflowScheduler) run(ctx context.Context) {
	var err error
	var nextSleep time.Duration
	for {
		select {
		case <-ctx.Done():
			logrus.Info("WorkflowScheduler shutting down")
			return
		case <-time.After(nextSleep):
			nextSleep, err = w.scheduleWorkflows(ctx)
			if err != nil {
				logrus.WithError(err).Error("failed to schedule workflows")
			}
		}
	}
}

func (w *workflowScheduler) scheduleWorkflows(ctx context.Context) (time.Duration, error) {
	toEnqueue, completer, err := w.backend.GetDueRecurringWorkflows(ctx)
	if err != nil {
		return maxWakeupInterval, errors.Wrap(err, "could not fetch recurring workflows")
	}
	defer completer.Cancel()

	sleepDuration := maxWakeupInterval
	for _, s := range toEnqueue {
		workflowInstanceName := fmt.Sprintf("%s/%s/%d", s.WorkflowName, s.Name, s.ID)

		// TODO(ssd) 2019-05-13: We might need two different
		// rule types here to suppor the different use cases.
		recurrence, err := rrule.StrToRRule(s.Recurrence)
		if err != nil {
			// TODO(ssd) 2019-05-13: Perhaps we should disable this rule so that it doesn't keep producing errors
			logrus.WithError(err).Error("could not parse recurrence rule for workflow, skipping")
			continue
		}

		nowUTC := time.Now().UTC()
		// NOTE(ssd) 2019-05-13: compliance looks 5 seconds in
		// the past to make sure that a job with a count of 1
		// actually gets run. However, I'm currently thinking
		// that we can push those kind of jobs onto the
		// workflow-instances queue immediately.
		nextDueAt := recurrence.After(nowUTC, true).UTC()
		err = completer.EnqueueRecurringWorkflow(s, workflowInstanceName, nextDueAt, nowUTC)
		if err != nil {
			if err == ErrWorkflowInstanceExists {
				// TODO(jaym): what do we want to do here? i think we're going to keep trying
				//             until we succeed here? Maybe we want to skip this interval?
				// It's also possible this happens on Commit instead of here
			}
			logrus.WithError(err).Error("could not update recurring workflow record")
			// TODO BUG (jaym): We cannot continue an errored transaction. We'll have to modify
			// the query to tell us when there is a conflict.
			// It seems the expected behavior when we try to commit this is to roll back.
			// What actually happens is we deadlock:
			// https://github.com/lib/pq/issues/731
			continue
		}

		if time.Until(nextDueAt) < sleepDuration {
			sleepDuration = time.Until(nextDueAt)
		}

	}

	return sleepDuration, completer.Commit()
}
