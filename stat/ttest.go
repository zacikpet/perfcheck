package stat

import "github.com/aclements/go-moremath/stats"

func makeSample(data []float64) stats.TTestSample {

	return stats.Sample{Xs: data}
}
