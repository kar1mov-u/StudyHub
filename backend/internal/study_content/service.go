package studycontent

import (
	"context"
	"io"
	"log/slog"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Queue interface {
	Consume() chan amqp.Delivery
}

type FileStorage interface {
	GetObject(ctx context.Context, objectID string) (io.ReadCloser, error)
}

type AI interface {
	Generate(ctx context.Context, file io.ReadCloser) (string, error)
}

type StudyContentRepository interface {
}

type StudyContentService struct {
	contentRepository StudyContentRepository
	queue             Queue
	delivery          chan amqp.Delivery
	fileStorage       FileStorage
	ai                AI
}

func NewStudyContentService(q Queue, fileStorage FileStorage, ai AI) *StudyContentService {
	s := &StudyContentService{
		queue:       q,
		fileStorage: fileStorage,
		ai:          ai,
		delivery:    q.Consume(),
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
	//workers will listen to the channel and consume the jobs
	for msg := range s.delivery {
		key := string(msg.Body)
		slog.Info("started on job with object id", "ID", string(msg.Body))

		//Download file to the disc
		file, err := s.fileStorage.GetObject(context.Background(), key)
		if err != nil {
			slog.Error("error getting file from storage", "err", err)
			continue
			//later may implement dead later queue
		}
		defer file.Close()

		result, err := s.ai.Generate(context.Background(), file)
		slog.Info("finished job ", " object id :", key, "result: ", result)
		//send the file to the LLM api

		//save to DB
		msg.Ack(true)
	}
	slog.Info("worker exiting")
}
