package handler

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTargetVerification(t *testing.T) {
	Convey("Target verification", t, func() {
		Convey("Check bad function", func() {
			target := `alias(test.one,'One'`
			request := targetVerification(target)
			So(request.SyntaxOk, ShouldBeFalse)
		})

		Convey("Check correct construction", func() {
			target := `alias(test.one,'One')`
			expected := targetVerification(target)
			So(expected.SyntaxOk, ShouldBeTrue)
		})

		for funcName, expected := range badFunctions.getCollection() {
			Convey("Check bad function: "+funcName, func() {
				target := funcName + `(seriesList, intervalString, func='sum', alignToFrom=False)`
				request := targetVerification(target)

				actual, ok := request.BadFunctions[funcName]

				So(request.SyntaxOk, ShouldBeTrue)
				So(expected.Type, ShouldEqual, actual.Type)
				So(ok, ShouldBeTrue)
			})
		}
	})
}
