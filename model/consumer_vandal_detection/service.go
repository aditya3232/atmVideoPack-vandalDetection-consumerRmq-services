package consumer_vandal_detection

type Service interface {
	ConsumerQueueVandalDetection() (RmqConsumerVandalDetection, error)
}

type service struct {
	vandalDetectionRepository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// consume and save to db
func (s *service) ConsumerQueueVandalDetection() (RmqConsumerVandalDetection, error) {

	// consume queue
	newRmqConsumerVandalDetection, err := s.vandalDetectionRepository.ConsumerQueueVandalDetection()
	if err != nil {
		return newRmqConsumerVandalDetection, err
	}

	return newRmqConsumerVandalDetection, nil

}
