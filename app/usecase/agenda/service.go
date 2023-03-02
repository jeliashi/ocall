package agenda

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"ocall/app/entity"
)

type Service struct {
	repo Repository
}

func NewService(repository Repository) Service {
	return Service{repository}
}

func (s *Service) CreateApplication(
	ctx context.Context, name string, status entity.ApplicationStatus,
	performer *entity.Performer, event *entity.Event, act *entity.Act,
) (entity.Application, error) {
	if status == entity.UNKNOWN {
		status = entity.PENDING
	}
	if performer == nil || performer.ID == uuid.Nil {
		_performer, err := entity.PerformerFromContext(ctx)
		if err != nil {
			return entity.Application{}, errors.Wrap(err, "unable to retrieve performer from context")
		}
		performer = &_performer
	}
	if event.ID == uuid.Nil {
		_event, err := entity.EventFromContext(ctx)
		if err != nil {
			return entity.Application{}, errors.Wrap(err, "unable to retrieve event from context")
		}
		event = &_event
	}
	a, err := entity.NewApplication(name, status, performer, event, act)
	if err != nil {
		return entity.Application{}, errors.Wrap(err, "unable to create application")
	}

	id, err := s.repo.CreateApplication(ctx, a)
	if err != nil {
		return entity.Application{}, errors.Wrap(err, "unable to save act")
	}
	a.ID = id
	return a, nil
}

func (s *Service) GetApplication(ctx context.Context, id uuid.UUID) (entity.Application, error) {
	a, err := s.repo.GetApplicationByID(ctx, id)
	if err != nil {
		return entity.Application{}, errors.Wrap(err, "unable to retrieve application")
	}
	return a, nil
}

func (s *Service) GetApplicationsByEvent(ctx context.Context, event *entity.Event) ([]entity.Application, error) {
	if event == nil || event.ID == uuid.Nil {
		_event, err := entity.EventFromContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to retrieve event by context")
		}
		event = &_event
	}
	applications, err := s.repo.GetApplicationByEventID(ctx, event.ID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get applications by event id")
	}
	return applications, nil
}

func (s *Service) GetApplicationsByPerformer(ctx context.Context, performer *entity.Performer) ([]entity.Application, error) {
	if performer == nil || performer.ID == uuid.Nil {
		_p, err := entity.PerformerFromContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get performers from context")
		}
		performer = &_p
	}
	applications, err := s.repo.GetApplicationsByPerformerID(ctx, performer.ID)
	if err != nil {
		return nil, errors.Wrap(err, "unable to query applications by performer id")
	}
	return applications, nil
}
func (s *Service) UpdateApplication(ctx context.Context, application entity.Application) error {
	if application.ID == uuid.Nil {
		return errors.New("application id cannot be 0")
	}
	err := s.repo.UpdateApplication(ctx, application)
	if err != nil {
		return errors.Wrap(err, "unable to update application")
	}
	return nil
}
func (s *Service) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	err := s.repo.DeleteApplication(ctx, id)
	if err != nil {
		return errors.Wrap(err, "unable to delete application")
	}
	return nil
}

func (s *Service) CreateEvent(
	ctx context.Context, name, description string,
	tags []entity.Tag, producer *entity.Producer,
	form entity.ApplicationForm,
) (entity.Event, error) {
	if producer == nil || producer.ID == uuid.Nil {
		_producer, err := entity.ProducerFromContext(ctx)
		if err != nil {
			return entity.Event{}, errors.Wrap(err, "unabe to get producer from context")
		}
		producer = &_producer
	}
	e, err := entity.NewEvent(name, description, tags, producer, form, nil, nil)
	if err != nil {
		return entity.Event{}, errors.Wrap(err, "unable to create event")
	}
	id, err := s.repo.CreateEvent(ctx, e)
	if err != nil {
		return entity.Event{}, errors.Wrap(err, "unable to save event")
	}
	e.ID = id
	return e, nil
}

func (s *Service) GetEvent(ctx context.Context, id uuid.UUID) (entity.Event, error) {
	e, err := s.repo.GetEventByID(ctx, id)
	if err != nil {
		return entity.Event{}, errors.Wrap(err, "unable to get event by id")
	}
	return e, nil
}
func (s *Service) AddPerformersToEvent(ctx context.Context, event *entity.Event, performers []entity.Performer) error {
	if event == nil || event.ID == uuid.Nil {
		_event, err := entity.EventFromContext(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to get event from context")
		}
		event = &_event
	}
	newActs := []entity.Application{}
	for _, performer := range performers {
		newAct, err := s.CreateApplication(
			ctx, fmt.Sprintf("Manual: %s - %s", performer.Name, event.Name),
			entity.PENDING, &performer, event, nil,
		)
		if err != nil {
			return errors.Wrapf(err, "unable to create application for performer %s", performer.Name)
		}
		newActs = append(newActs, newAct)
	}
	err := s.repo.AddApplicationsToEvent(ctx, *event, performers)
	if err != nil {
		return errors.Wrap(err, "unable to save performers to event")
	}
	return nil
}

func (s *Service) GetEventsByProducer(ctx context.Context, producer *entity.Producer) ([]entity.Event, error) {
	if producer == nil || producer.ID == uuid.Nil {
		_producer, err := entity.ProducerFromContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get producer from context")
		}
		producer = &_producer
	}
	events, err := s.repo.GetEventsByProducerID(ctx, producer.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get events by producer id %s", producer.ID)
	}
	return events, nil
}

func (s *Service) GetEventsByVenue(ctx context.Context, venue *entity.Venue) ([]entity.Event, error) {
	if venue == nil || venue.ID == uuid.Nil {
		_venue, err := entity.VenueFromContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get venue by context")
		}
		venue = &_venue
	}
	events, err := s.repo.GetEventsByVenueID(ctx, venue.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get events for venue id %s", venue.ID)
	}
	return events, nil
}

func (s *Service) UpdateEvent(ctx context.Context, event entity.Event) error {
	if event.ID == uuid.Nil {
		return errors.New("event id must be 0")
	}
	err := s.repo.UpdateEvent(ctx, event)
	if err != nil {
		return errors.Wrap(err, "unable to commit changes")
	}
	return nil
}

func (s *Service) ChangeEventApplicationStatus(ctx context.Context, id uuid.UUID, status entity.EventApplicationStatus) error {
	err := s.repo.ChangeEventApplicationStatus(ctx, id, status)
	if err != nil {
		return errors.Wrap(err, "unable to update event application status")
	}
	return nil
}

func (s *Service) ModifyActListForEvent(ctx context.Context, event *entity.Event, acts []entity.Act) error {
	if event == nil || event.ID == uuid.Nil {
		_event, err := entity.EventFromContext(ctx)
		if err != nil {
			return errors.Wrap(err, "unable to get event from context")
		}
		event = &_event
	}
	if err := s.repo.ModifyActListForEvent(ctx, event.ID, acts); err != nil {
		return errors.Wrap(err, "unable to modify act list")
	}
	return nil
}

func (s *Service) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteEvent(ctx, id); err != nil {
		return errors.Wrap(err, "unable to delete event")
	}
	return nil
}

func (s *Service) CreateAct(
	ctx context.Context, name, description string,
	tags []entity.Tag, media []entity.Media,
	form entity.ApplicationForm, performer *entity.Performer,
) (entity.Act, error) {
	if performer == nil || performer.ID == uuid.Nil {
		_performer, err := entity.PerformerFromContext(ctx)
		if err != nil {
			return entity.Act{}, errors.Wrap(err, "unable to get performer from context")
		}
		performer = &_performer
	}
	a, err := entity.NewAct(name, tags, description, media, form, performer)
	if err != nil {
		return entity.Act{}, errors.Wrap(err, "unable to create new act")
	}
	id, err := s.repo.CreateAct(ctx, a)
	if err != nil {
		return entity.Act{}, errors.Wrap(err, "unable to commit act")
	}
	a.ID = id
	return a, nil
}

func (s *Service) GetActs(ctx context.Context, ids ...uuid.UUID) ([]entity.Act, error) {
	acts, err := s.repo.GetActs(ctx, ids)
	if err != nil {
		return nil, errors.Wrap(err, "unable to retrieve acts by id")
	}
	return acts, nil
}

func (s *Service) GetActsByPerformer(ctx context.Context, performer *entity.Performer) ([]entity.Act, error) {
	if performer == nil || performer.ID == uuid.Nil {
		_performer, err := entity.PerformerFromContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get performer from context")
		}
		performer = &_performer
	}
	acts, err := s.repo.GetActsByPerformerID(ctx, performer.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get acts by performer id %s", performer.ID)
	}
	return acts, nil
}

func (s *Service) GetConfirmedActsByEvent(ctx context.Context, event *entity.Event) ([]entity.Act, error) {
	if event == nil || event.ID == uuid.Nil {
		_event, err := entity.EventFromContext(ctx)
		if err != nil {
			return nil, errors.New("unable to get event by context")
		}
		event = &_event
	}
	acts, err := s.repo.GetConfirmedActsByEventID(ctx, event.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to query acts by event id %s", event.ID)
	}
	return acts, nil
}

func (s *Service) GetSavedActByProducer(ctx context.Context, producer *entity.Producer) ([]entity.Act, error) {
	if producer == nil || producer.ID == uuid.Nil {
		_producer, err := entity.ProducerFromContext(ctx)
		if err != nil {
			return nil, errors.Wrap(err, "unable to get producer from context")
		}
		producer = &_producer
	}
	acts, err := s.repo.GetSavedActsByProducerID(ctx, producer.ID)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to get saved acts by producer id %s", producer.ID)
	}
	return acts, nil
}

func (s *Service) UpdateAct(ctx context.Context, act entity.Act) error {
	if act.ID == uuid.Nil {
		return errors.New("act id must not be 0")
	}
	err := s.repo.UpdateAct(ctx, act)
	if err != nil {
		return errors.Wrap(err, "unable to commit changes")
	}
	return nil
}

func (s *Service) DeleteAct(ctx context.Context, id uuid.UUID) error {
	if err := s.repo.DeleteAct(ctx, id); err != nil {
		return errors.Wrap(err, "unable to delete act")
	}
	return nil
}
