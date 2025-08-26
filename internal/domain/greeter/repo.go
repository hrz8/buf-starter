package greeter

import "fmt"

var allowedNames = []string{
	"Alina", "Bryce", "Carmen", "Darius", "Elena",
	"Felix", "Gianna", "Hassan", "Irene", "Jasper",
	"Kiana", "Luther", "Maya", "Nolan", "Orlando",
	"Priya", "Quincy", "Rafael", "Sienna", "Tobias",
	"Umair", "Vera", "Wesley", "Xavier", "Yasmin",
	"Zane", "Adriana", "Bennett", "Clarissa", "Devonte",
	"Estella", "Finnegan", "Gracelyn", "Harvey", "Isidora",
	"Jovani", "Katarina", "Leonidas", "Mirella", "Nikolas",
	"Octavia", "Percival", "Quintessa", "Romero", "Salvador",
	"Theodora", "Ulrich", "Valeria", "Winslow", "Xiomara",
	"Yuridia", "Zephyrus", "Aurelius", "Bellatrix", "Caspian",
	"Demetrius", "Evangeline", "Florentino", "Galadriel", "Hermione",
	"Ignatius", "Julianna", "Kristoffer", "Lysandra", "Maximiliano",
	"Nefertari", "Olivander", "Philomena", "Quetzalcoatl", "Rhiannon",
	"Sebastiana", "Thessalonia", "Ulyssiana", "Vladimir", "Wilhelmina",
	"Xenophilius", "Yggdrasila", "Zaphkiel", "Alejandrina", "Balthazar",
	"Christabelle", "Domenico", "Euphrosyne", "Featherstone", "Gwendolyn",
	"Hyacinthus", "Isambard", "Jacqueline", "Kallistrate", "Leontius",
	"Marcellinus", "Nicomachus", "Ozymandias", "Petronella", "Quintilius",
	"Rosencrantz", "Seraphimiel", "Timotheus", "Ultraviolet", "Valentinian",
}

type Repo struct{}

func NewRepo() *Repo {
	return &Repo{}
}

func (r *Repo) GetGreeterTemplate(name string) string {
	return fmt.Sprintf("Hello, %s!", name)
}

func (r *Repo) GetAllowedNames(page, limit int32) []string {
	start := int((page - 1) * limit)
	if start > len(allowedNames) {
		return []string{}
	}
	end := min(start+int(limit), len(allowedNames))
	return allowedNames[start:end]
}

// New method that returns both paginated names and total count
func (r *Repo) GetAllowedNamesWithTotal(page, limit int32) ([]string, int32) {
	total := int32(len(allowedNames))
	names := r.GetAllowedNames(page, limit)
	return names, total
}

// New method to get total count only
func (r *Repo) GetTotalAllowedNames() int32 {
	return int32(len(allowedNames))
}
