package consumer_vandal_detection

type Service interface {
	ConsumerQueueVandalDetection() (RmqConsumerVandalDetection, error)
}

type service struct {
	statusMcDetectionRepository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// consume and save to db
func (s *service) ConsumerQueueVandalDetection() (RmqConsumerVandalDetection, error) {

	// consume queue
	newRmqConsumerVandalDetection, err := s.statusMcDetectionRepository.ConsumerQueueVandalDetection()
	if err != nil {
		return newRmqConsumerVandalDetection, err
	}

	return newRmqConsumerVandalDetection, nil

}
