package util

import (
	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"testing"
)

func TestLogger(t *testing.T) {

	Convey("GetLogger", t, func() {

		Convey("shall return non nil instance", func() {
			So(GetLogger(zapcore.InfoLevel), ShouldNotBeNil)
		})

		Convey("shall respect log level", func() {
			logger := GetLogger(zapcore.InfoLevel).(*zap.SugaredLogger)
			So(logger.Desugar().Core().Enabled(zap.DebugLevel), ShouldBeFalse)
			So(logger.Desugar().Core().Enabled(zap.InfoLevel), ShouldBeTrue)
		})

		Convey("shall just work", func() {
			logger := GetLogger(zapcore.InfoLevel).(*zap.SugaredLogger)
			logger.Info("")
		})

		Convey("shall panic when log level is not valid ", func() {
			So(func() {
				logger := GetLogger(zapcore.Level(123))
				logger.Info("Should panic earlier")
			}, ShouldPanic)
		})

	})
}
