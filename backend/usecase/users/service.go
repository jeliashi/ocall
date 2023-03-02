package users

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"ocall/app/entity"
)

type Service struct {
	repo Repository
}

func NewService(repository Repository) (Service, error) {
	return Service{repository}, nil
}

func (s *Service) CreatePerformer(ctx context.Context, name string, userID uuid.UUID) (entity.Performer, error) {
	p, err := entity.NewPerformer(name, userID)
	if err != nil {
		return entity.Performer{}, errors.Wrap(err, "unable to create performer")
	}
	var id uuid.UUID
	id, err = s.repo.CreatePerformer(ctx, p)
	if err != nil {
		return entity.Performer{}, errors.Wrap(err, "unable to commit performer")
	}
	p.ID = id
	return p, nil
}

func (s *Service) GetPerformerByID(ctx context.Context, id uuid.UUID) (entity.Performer, error) {
	p, err := s.repo.GetPerformerByID(ctx, id)
	if err != nil {
		return entity.Performer{}, errors.Wrap(err, "unable to retrieve")
	}
	return p, nil
}

func (s *Service) GetPerformersByProducer(ctx context.Context, producer *entity.Producer) ([]entity.Performer, error) {
	var err error
	if producer == nil || producer.ID == uuid.Nil {
		_producer, err := entity.ProducerFromContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get producer from context")
		}
		producer = &_producer
	}
	performers, err := s.repo.GetPerformersByProducerID(ctx, producer.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get performers by producer id %s", producer.ID)
	}
	return performers, nil
}

func (s *Service) GetPerformersByEvent(ctx context.Context, event *entity.Event) ([]entity.Performer, error) {
	var err error
	if event == nil || event.ID == uuid.Nil {
		_event, err := entity.EventFromContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get event from context")
		}
		event = &_event
	}
	performers, err := s.repo.GetPerformersByEventID(ctx, event.ID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get performers by event id")
	}
	return performers, nil
}

func (s *Service) UpdatePerformer(ctx context.Context, performer entity.Performer) error {
	err := s.repo.UpdatePerformer(ctx, performer)
	if err != nil {
		return errors.Wrap(err, "user app:")
	}
	return nil
}

func (s *Service) DeletePerformer(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeletePerformer(ctx, id)
	if err != nil {
		return errors.Wrap(err, "user app:")
	}
	return nil
}

func (s *Service) CreateProducer(ctx context.Context, name string, userID uuid.UUID) (entity.Producer, error) {
	p, err := entity.NewProducer(name, userID)
	if err != nil {
		return entity.Producer{}, errors.Wrap(err, "user app:")
	}
	id, err := s.repo.CreateProducer(ctx, p)
	if err != nil {
		return entity.Producer{}, errors.Wrap(err, "user app unable to commit:")
	}
	p.ID = id
	return p, nil
}

func (s *Service) GetProducer(ctx context.Context, id uuid.UUID) (entity.Producer, error) {
	p, err := s.repo.GetProducerByID(ctx, id)
	if err != nil {
		return entity.Producer{}, errors.Wrap(err, "user app")
	}
	return p, nil
}

func (s *Service) GetProducerByEvent(ctx context.Context, event entity.Event) (entity.Producer, error) {
	var err error
	if event.ID == uuid.Nil {
		event, err = entity.EventFromContext(ctx)
		if err != nil {
			return entity.Producer{}, errors.Wrap(err, "unable to retrieve event from context")
		}
	}
	producer, err := s.repo.GetProducerByEventID(ctx, event.ID)
	if err != nil {
		return entity.Producer{}, errors.Wrap(err, "user service")
	}
	return producer, nil
}

func (s *Service) UpdateProducer(ctx context.Context, producer entity.Producer) error {
	err := s.repo.UpdateProducer(ctx, producer)
	if err != nil {
		return errors.Wrap(err, "user app")
	}
	return nil
}
func (s *Service) DeleteProducer(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteProducer(ctx, id)
	if err != nil {
		return errors.Wrap(err, "user app")
	}
	return nil
}

func (s *Service) CreateVenue(ctx context.Context, name string, userID uuid.UUID) (entity.Venue, error) {
	v, err := entity.NewVenue(name, userID)
	if err != nil {
		return entity.Venue{}, errors.Wrap(err, "user app")
	}
	id, err := s.repo.CreateVenue(ctx, v)
	if err != nil {
		return entity.Venue{}, errors.Wrap(err, "user app")
	}
	v.ID = id
	return v, nil
}

func (s *Service) GetVenue(ctx context.Context, id uuid.UUID) (entity.Venue, error) {
	if id == uuid.Nil {
		return entity.Venue{}, errors.New("invalid nil UUID")
	}
	v, err := s.repo.GetVenueByID(ctx, id)
	if err != nil {
		return entity.Venue{}, errors.Wrap(err, "user app")
	}
	return v, nil
}

func (s *Service) GetVenueByEvent(ctx context.Context, event *entity.Event) (entity.Venue, error) {
	if event == nil || event.ID == uuid.Nil {
		_event, err := entity.EventFromContext(ctx)
		if err != nil {
			return entity.Venue{}, errors.Wrap(err, "unable to get event from context")
		}
		event = &_event
	}
	venue, err := s.repo.GetVenueByEventID(ctx, event.ID)
	if err != nil {
		return entity.Venue{}, errors.Wrap(err, "unable to get venue by event id")
	}
	return venue, nil
}

func (s *Service) UpdateVenue(ctx context.Context, venue entity.Venue) error {
	err := s.repo.UpdateVenue(ctx, venue)
	if err != nil {
		return errors.Wrap(err, "user app")
	}
	return nil
}

func (s *Service) DeleteVenue(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteVenue(ctx, id)
	if err != nil {
		return errors.Wrap(err, "user ser ice")
	}
	return nil
}

func (s *Service) ListProfiles(ctx context.Context, userID uuid.UUID) ([]entity.Performer, []entity.Producer, []entity.Venue, error) {
	performers, err := s.repo.GetPerformersByUser(ctx, userID)
	if err != nil {
		return nil, nil, nil, errors.Wrap(err, "user app unable to query performers")
	}
	producers, err := s.repo.GetProducersByUser(ctx, userID)
	if err != nil {
		return performers, nil, nil, errors.Wrap(err, "user app unable to query producers")
	}
	venues, err := s.repo.GetVenuesByUser(ctx, userID)
	if err != nil {
		return performers, producers, nil, errors.Wrap(err, "user app unable to query producers")
	}
	return performers, producers, venues, nil
}
