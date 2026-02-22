package studycontent

import (
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Consume() chan amqp.Delivery
}

type StudyContentService struct {
	queue    Queue
	delivery chan amqp.Delivery
}

func NewStudyContentService(q Queue) *StudyContentService {
	s := &StudyContentService{
		queue:    q,
		delivery: q.Consume(),
	}
	s.startWorkers()
	return s
}

func (s *StudyContentService) startWorkers() {
	for i := range 5 {
		slog.Info("started worker", "id", i)
		go s.worker()
	}
}

func (s *StudyContentService) worker() {
	for msg := range s.delivery {
		slog.Info("started on job with object id", "ID", string(msg.Body))
		//get the file from the s3

		//send the file to the ai

		//save to DB
		msg.Ack(true)
	}
	slog.Info("worker exiting")
}
