package graph

import (
	"digi-bot/model"
	"fmt"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"math/rand"
	"os"
	"time"
)

func LinearRegreasion(histories []model.History) (string, error) {
	var xvalues []time.Time
	var yvalues []float64
	for _, history := range histories {
		xvalues = append(xvalues, history.Date)
		yvalues = append(yvalues, float64(history.Price))
	}

	mainSeries := chart.TimeSeries{
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0).WithAlpha(64),
			FillColor:   chart.GetDefaultColor(0).WithAlpha(64),
			Padding: chart.Box{
				Top:    50,
				Left:   25,
				Right:  25,
				Bottom: 10,
			},
		},
		Name:    "A test series",
		XValues: xvalues,
		YValues: yvalues,
	}

	linRegSeries := &chart.LinearRegressionSeries{
		InnerSeries: mainSeries,
	} // we can optionally set the `WindowSize` property which alters how the moving average is calculated.

	graph := chart.Chart{
		Series: []chart.Series{
			mainSeries,
			linRegSeries,
		},
	}

	graph.Elements = []chart.Renderable{
		chart.Legend(&graph),
	}

	filename := fmt.Sprintf("%d.png", rand.Int())
	f, _ := os.Create(filename)
	defer f.Close()
	err := graph.Render(chart.PNG, f)

	return filename, err
}

func StockAnalysis(histories []model.History) string {
	var xvalues []time.Time
	var yvalues []float64
	for _, history := range histories {
		xvalues = append(xvalues, history.Date)
		yvalues = append(yvalues, float64(history.Price))
	}

	priceSeries := chart.TimeSeries{
		Name: "SPY",
		Style: chart.Style{
			StrokeColor: chart.GetDefaultColor(0),
		},
		XValues: xvalues,
		YValues: yvalues,
	}

	smaSeries := chart.SMASeries{
		Name: "SPY - SMA",
		Style: chart.Style{
			StrokeColor:     drawing.ColorRed,
			StrokeDashArray: []float64{5.0, 5.0},
		},
		InnerSeries: priceSeries,
	}

	bbSeries := &chart.BollingerBandsSeries{
		Name: "SPY - Bol. Bands",
		Style: chart.Style{
			StrokeColor: drawing.ColorFromHex("efefef"),
			FillColor:   drawing.ColorFromHex("efefef").WithAlpha(64),
		},
		InnerSeries: priceSeries,
	}

	graph := chart.Chart{
		XAxis: chart.XAxis{
			TickPosition: chart.TickPositionBetweenTicks,
		},
		YAxis: chart.YAxis{
			Range: &chart.ContinuousRange{
				Max: 220.0,
				Min: 180.0,
			},
		},
		Series: []chart.Series{
			bbSeries,
			priceSeries,
			smaSeries,
		},
	}

	f, _ := os.Create("STOCK.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
	return ""
}