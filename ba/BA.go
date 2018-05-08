package ba

import (
	"../graph"
	"../jssp"
	"../constants"
	"../gantt"
	"fmt"
	"sort"
	"math/rand"
	"strconv"
)

func LocalSearch(solution jssp.Solution) jssp.Solution {
	neighbours := jssp.GenerateNeighbourHood(solution)
	return neighbours[rand.Intn(len(neighbours))]
}


func BA(problemGraph graph.Graph) {
	fmt.Println("Running BA")
	nBees := (constants.NElites * constants.NRE + (constants.NBest - constants.NElites) * constants.NRB) + constants.Scouts
	fmt.Println("Number of bees:", nBees)

	bees := make([]jssp.Solution, 0)

	for i := 0; i < constants.Scouts; i++ {
		bees = append(bees, jssp.ListScheduler(problemGraph, nil, true))
	}

	sort.Slice(bees, func(i, j int) bool {
		return jssp.MakeSpan2(bees[i]) < jssp.MakeSpan2(bees[j])
	})

	for c := 0; c < 10; c++{
		fmt.Println("Iteration", c, "Best solution", jssp.MakeSpan2(bees[0]))
		// Recruit
		sort.Slice(bees, func(i, j int) bool {
			return jssp.MakeSpan2(bees[i]) < jssp.MakeSpan2(bees[j])
		})
		bees = bees[0:constants.NBest]
		for i := 0; i < constants.NElites; i++ {
			for j := 0; j < constants.NRE; j++ {
				bees = append(bees, LocalSearch(bees[i]))
			}
		}
		for i := constants.NElites; i < constants.NBest; i++ {
			for j := 0; j < constants.NRB; j++ {
				bees = append(bees, LocalSearch(bees[i]))
			}
		}
		for i := len(bees); i < nBees; i++ {
			bees = append(bees, jssp.ListScheduler(problemGraph, nil, true))
		}
	}

	fmt.Println("Done!")
	fmt.Println("Makespan:", jssp.CalculateMakespan(bees[0]))
	orders := graph.NodeListToOrderList(bees[0].Nodes, bees[0].StartTimeMap)
	gantt.CreateChart("03 - Program Outputs/BEE_Chart_" + strconv.Itoa(constants.ProblemNumber) + ".xlsx", orders)


}