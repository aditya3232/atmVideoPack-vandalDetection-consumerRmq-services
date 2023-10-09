package consumer_vandal_detection

type Service interface {
	ConsumerQueueVandalDetection() (VandalDetection, error)
}

type service struct {
	vandalDetectionRepository Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

// consume and save to db
func (s *service) ConsumerQueueVandalDetection() (VandalDetection, error) {

	// consume queue
	vandalDetection, err := s.vandalDetectionRepository.ConsumerQueueVandalDetection()
	if err != nil {
		return VandalDetection{}, err
	}

	return vandalDetection, nil

}
