package greeter

import (
	"context"

	"buf.build/go/protovalidate"
	"github.com/hrz8/altalune"

	greeterv1 "github.com/hrz8/altalune/gen/greeter/v1"
)

type Service struct {
	greeterv1.UnimplementedGreeterServiceServer

	validator   protovalidate.Validator
	log         altalune.Logger
	greeterRepo Repositor
}

func NewService(v protovalidate.Validator, log altalune.Logger, greeterRepo Repositor) *Service {
	return &Service{
		validator:   v,
		log:         log,
		greeterRepo: greeterRepo,
	}
}

var allowedNameMap = map[string]bool{
	"Alina": true, "Bryce": true, "Carmen": true, "Darius": true, "Elena": true,
	"Felix": true, "Gianna": true, "Hassan": true, "Irene": true, "Jasper": true,
	"Kiana": true, "Luther": true, "Maya": true, "Nolan": true, "Orlando": true,
	"Priya": true, "Quincy": true, "Rafael": true, "Sienna": true, "Tobias": true,
	"Umair": true, "Vera": true, "Wesley": true, "Xavier": true, "Yasmin": true,
	"Zane": true, "Adriana": true, "Bennett": true, "Clarissa": true, "Devonte": true,
	"Estella": true, "Finnegan": true, "Gracelyn": true, "Harvey": true, "Isidora": true,
	"Jovani": true, "Katarina": true, "Leonidas": true, "Mirella": true, "Nikolas": true,
	"Octavia": true, "Percival": true, "Quintessa": true, "Romero": true, "Salvador": true,
	"Theodora": true, "Ulrich": true, "Valeria": true, "Winslow": true, "Xiomara": true,
	"Yuridia": true, "Zephyrus": true, "Aurelius": true, "Bellatrix": true, "Caspian": true,
	"Demetrius": true, "Evangeline": true, "Florentino": true, "Galadriel": true, "Hermione": true,
	"Ignatius": true, "Julianna": true, "Kristoffer": true, "Lysandra": true, "Maximiliano": true,
	"Nefertari": true, "Olivander": true, "Philomena": true, "Quetzalcoatl": true, "Rhiannon": true,
	"Sebastiana": true, "Thessalonia": true, "Ulyssiana": true, "Vladimir": true, "Wilhelmina": true,
	"Xenophilius": true, "Yggdrasila": true, "Zaphkiel": true, "Alejandrina": true, "Balthazar": true,
	"Christabelle": true, "Domenico": true, "Euphrosyne": true, "Featherstone": true, "Gwendolyn": true,
	"Hyacinthus": true, "Isambard": true, "Jacqueline": true, "Kallistrate": true, "Leontius": true,
	"Marcellinus": true, "Nicomachus": true, "Ozymandias": true, "Petronella": true, "Quintilius": true,
	"Rosencrantz": true, "Seraphimiel": true, "Timotheus": true, "Ultraviolet": true, "Valentinian": true,
}

func (s *Service) SayHello(ctx context.Context, req *greeterv1.SayHelloRequest) (*greeterv1.SayHelloResponse, error) {
	if err := s.validator.Validate(req); err != nil {
		return nil, altalune.NewInvalidPayloadError(err.Error())
	}
	if _, ok := allowedNameMap[req.Name]; !ok {
		return nil, altalune.NewGreetingUnrecognize(req.Name)
	}
	msg := s.greeterRepo.GetGreeterTemplate(req.Name)
	response := &greeterv1.SayHelloResponse{
		Message: msg,
	}
	return response, nil
}
