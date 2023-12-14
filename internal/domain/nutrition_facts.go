//go:generate mockgen -destination=./mocks/nutrition_facts.go -package=mocks -source=nutrition_facts.go

package domain

type NutritionFacts struct {
	Calories int
	Fat      float64
	Carbs    float64
	Protein  float64
	Fiber    float64
}

func (nutritionFacts NutritionFacts) ToDto() NutritionFactsDto {
	return NutritionFactsDto{
		Calories: nutritionFacts.Calories,
		Fat:      nutritionFacts.Fat,
		Carbs:    nutritionFacts.Carbs,
		Protein:  nutritionFacts.Protein,
		Fiber:    nutritionFacts.Fiber,
	}
}

type NutritionFactsDto struct {
	Calories int     `json:"calories"`
	Fat      float64 `json:"fat"`
	Carbs    float64 `json:"carbs"`
	Protein  float64 `json:"protein"`
	Fiber    float64 `json:"fiber"`
}

func (nutritionFactsDto NutritionFactsDto) ToEntity() NutritionFacts {
	return NutritionFacts{
		Calories: nutritionFactsDto.Calories,
		Fat:      nutritionFactsDto.Fat,
		Carbs:    nutritionFactsDto.Carbs,
		Protein:  nutritionFactsDto.Protein,
		Fiber:    nutritionFactsDto.Fiber,
	}
}
