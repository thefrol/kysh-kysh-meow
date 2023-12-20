package storage_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app"
	"github.com/thefrol/kysh-kysh-meow/internal/server/app/manager"
	"github.com/thefrol/kysh-kysh-meow/internal/server/storage"
)

type AdapterSuite struct {
	suite.Suite
	counters manager.CounterRepository
	gauges   manager.GaugeRepository
}

func (suite *AdapterSuite) SetupTest() {
	s := storage.New()
	operator := storage.AsOperator(s)

	suite.gauges = &storage.GaugeAdapter{Op: operator}
	suite.counters = &storage.CounterAdapter{Op: operator}
}

func (suite *AdapterSuite) TestGauges() {
	ctx := context.Background()
	suite.Run("not found", func() {
		_, err := suite.gauges.Gauge(ctx, "some_id")
		suite.ErrorIs(err, app.ErrorMetricNotFound)
	})

	suite.Run("Set&Get", func() {
		const (
			id   = "some_id"
			val1 = 1.01
			val2 = 1.02
		)

		v, err := suite.gauges.GaugeUpdate(ctx, id, val1)
		suite.NoError(err)
		suite.Equal(val1, v)

		v, err = suite.gauges.Gauge(ctx, id)
		suite.NoError(err)
		suite.Equal(val1, v)

		v, err = suite.gauges.GaugeUpdate(ctx, id, val2)
		suite.NoError(err)
		suite.Equal(val2, v)

		v, err = suite.gauges.Gauge(ctx, id)
		suite.NoError(err)
		suite.Equal(val2, v)

	})
}

func (suite *AdapterSuite) TestCounters() {
	ctx := context.Background()
	suite.Run("not found", func() {
		_, err := suite.counters.Counter(ctx, "some_id")
		suite.ErrorIs(err, app.ErrorMetricNotFound)
	})

	suite.Run("Set&Get", func() {
		const (
			id   = "some_id"
			val1 = 1
			val2 = 2
		)

		v, err := suite.counters.CounterIncrement(ctx, id, val1)
		suite.NoError(err)
		suite.EqualValues(val1, int64(v))

		v, err = suite.counters.Counter(ctx, id)
		suite.NoError(err)
		suite.EqualValues(val1, int64(v))

		v, err = suite.counters.CounterIncrement(ctx, id, val2)
		suite.NoError(err)
		suite.EqualValues(val2+val1, int64(v))

		v, err = suite.counters.Counter(ctx, id)
		suite.NoError(err)
		suite.EqualValues(val2+val1, int64(v))

	})
}

func Test_AppAdapters(t *testing.T) {
	suite.Run(t, new(AdapterSuite))
}
