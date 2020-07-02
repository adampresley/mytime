package sessions

import (
	"fmt"
	"time"

	"github.com/adampresley/mytime/api/categories"
	"github.com/adampresley/mytime/api/clients"
	"github.com/adampresley/mytime/api/helpers"
	"github.com/adampresley/mytime/api/projects"
	"github.com/adampresley/simdb"
)

type SessionServicer interface {
	CloseSession(sessionID int) error
	CreateSession(session Session) (string, error)
	DeleteActiveSessions() error
	HasActiveSession() (bool, error)
	GetActiveSession() (ActiveSession, error)
	GetSessionByID(sessionID int) (Session, error)
	InvoiceSessions(sessionIDs []int) []error
	InvoiceSession(sessionID int) error
	ListSessions(search SessionSearch) (SessionCollection, error)
	StartActiveSession(projectID, categoryID, clientID int, notes string) (ActiveSession, time.Time, error)
}

type SessionServiceConfig struct {
	CategoryService categories.CategoryServicer
	ClientService   clients.ClientServicer
	DB              *simdb.Driver
	HelperService   helpers.HelperServicer
	ProjectService  projects.ProjectServicer
}

type SessionService struct {
	CategoryService categories.CategoryServicer
	ClientService   clients.ClientServicer
	DB              *simdb.Driver
	HelperService   helpers.HelperServicer
	ProjectService  projects.ProjectServicer
}

func NewSessionService(config SessionServiceConfig) SessionService {
	return SessionService{
		CategoryService: config.CategoryService,
		ClientService:   config.ClientService,
		DB:              config.DB,
		HelperService:   config.HelperService,
		ProjectService:  config.ProjectService,
	}
}

func (s SessionService) CloseSession(sessionID int) error {
	var err error
	var session Session

	if session, err = s.GetSessionByID(sessionID); err != nil {
		return fmt.Errorf("Cannot find session %d: %w", sessionID, err)
	}

	if session.Invoiced && session.Paid {
		return fmt.Errorf("Session already closed!")
	}

	session.Invoiced = true
	session.InvoiceDate = time.Now()
	session.Paid = true
	session.PaidDate = time.Now()

	err = s.DB.Open(Session{}).Update(session)

	if err != nil {
		return fmt.Errorf("Error updating session %d: %w", sessionID, err)
	}

	return nil
}

func (s SessionService) CreateSession(session Session) (int, error) {
	session.SessionID = s.DB.Open(Session{}).GetNextNumericID()

	return session.SessionID, s.DB.Open(Session{}).Insert(session)
}

func (s SessionService) DeleteActiveSessions() error {
	var err error
	var activeSessions []ActiveSession

	if err = s.DB.Open(ActiveSession{}).Get().AsEntity(&activeSessions); err != nil {
		return fmt.Errorf("Error getting a list of active sessions: %w", err)
	}

	for _, as := range activeSessions {
		if err = s.DB.Open(ActiveSession{}).Delete(as); err != nil {
			return fmt.Errorf("Error deleting active session %s: %w", as.ActiveSessionID, err)
		}
	}

	return nil
}

func (s SessionService) GetActiveSession() (ActiveSession, error) {
	var (
		err              error
		hasActiveSession bool
		activeSessions   []ActiveSession
	)

	if hasActiveSession, err = s.HasActiveSession(); err != nil {
		return ActiveSession{}, err
	}

	if !hasActiveSession {
		return ActiveSession{}, fmt.Errorf("There is no active session")
	}

	if err = s.DB.Open(ActiveSession{}).Get().AsEntity(&activeSessions); err != nil {
		return ActiveSession{}, err
	}

	if len(activeSessions) < 1 {
		return ActiveSession{}, fmt.Errorf("There is no active sessions")
	}

	return activeSessions[0], nil
}

func (s SessionService) GetSessionByID(sessionID int) (Session, error) {
	var err error
	var session Session

	err = s.DB.Open(Session{}).Where("sessionID", "=", sessionID).First().AsEntity(&session)
	return session, err
}

func (s SessionService) HasActiveSession() (bool, error) {
	var err error
	var sessions []ActiveSession

	if err = s.DB.Open(ActiveSession{}).Get().AsEntity(&sessions); err != nil {
		return false, err
	}

	return len(sessions) > 0, nil
}

func (s SessionService) InvoiceSessions(sessionIDs []int) []error {
	result := make([]error, len(sessionIDs))

	for index, sessionID := range sessionIDs {
		err := s.InvoiceSession(sessionID)

		if err != nil {
			result[index] = fmt.Errorf("Session ID: %d - %w", sessionID, err)
		} else {
			result[index] = err
		}
	}

	return result
}

func (s SessionService) InvoiceSession(sessionID int) error {
	var err error
	var session Session

	if session, err = s.GetSessionByID(sessionID); err != nil {
		return err
	}

	if session.Invoiced {
		return fmt.Errorf("Session already invoiced")
	}

	if session.Paid {
		return fmt.Errorf("Cannot invoice a session that is already paid")
	}

	session.Invoiced = true
	session.InvoiceDate = time.Now()

	return s.UpdateSession(session)
}

func (s SessionService) ListSessions(search SessionSearch) (SessionCollection, error) {
	var err error

	// allSessions := make(SessionCollection, 0, 100)
	result := make(SessionCollection, 0, 100)

	d := s.DB.Open(Session{})

	if search.Paid {
		d = d.Where("paid", "=", true)
	} else {
		d = d.Where("paid", "=", false)
	}

	if search.Invoiced {
		d = d.Where("invoiced", "=", true)
	} else {
		d = d.Where("invoiced", "=", false)
	}

	if err = d.Get().AsEntity(&result); err != nil {
		return result, fmt.Errorf("Error querying for sessions: %w", err)
	}

	/*
	 * Start filtering
	 */
	filter := func(f func(session Session) bool) SessionCollection {
		r := make(SessionCollection, 0, 100)

		for _, s := range result {
			if f(s) {
				r = append(r, s)
			}
		}

		return r
	}

	if search.CategoryCode != "" {
		result = filter(func(session Session) bool {
			var category categories.Category

			if category, err = s.CategoryService.GetCategoryByID(session.CategoryID); err != nil {
				return false
			}

			if category.Code == search.CategoryCode {
				return true
			}

			return false
		})
	}

	if search.ClientCode != "" {
		result = filter(func(session Session) bool {
			var client clients.Client

			if client, err = s.ClientService.GetClientByID(session.ClientID); err != nil {
				return false
			}

			if client.Code == search.ClientCode {
				return true
			}

			return false
		})
	}

	if search.ProjectCode != "" {
		result = filter(func(session Session) bool {
			var project projects.Project

			if project, err = s.ProjectService.GetProjectByID(session.ProjectID); err != nil {
				return false
			}

			if project.Code == search.ProjectCode {
				return true
			}

			return false
		})
	}

	if search.SessionID > 0 {
		result = filter(func(session Session) bool {
			return session.SessionID == search.SessionID
		})
	}

	if len(search.SessionIDs) > 0 {
		result = filter(func(session Session) bool {
			result := false

			for _, id := range search.SessionIDs {
				if session.SessionID == id {
					result = true
					break
				}
			}

			return result
		})
	}

	return result, err
}

func (s SessionService) StartActiveSession(projectID, categoryID, clientID int, notes string) (ActiveSession, time.Time, error) {
	var err error
	var id string

	id = s.HelperService.GenerateID()
	startTime := time.Now()

	session := ActiveSession{
		ActiveSessionID: id,
		ProjectID:       projectID,
		ClientID:        clientID,
		CategoryID:      categoryID,
		StartTime:       startTime,
		Notes:           notes,
	}

	if err = s.DB.Open(ActiveSession{}).Insert(session); err != nil {
		return session, time.Now(), err
	}

	return session, startTime, nil
}

func (s SessionService) UpdateSession(session Session) error {
	return s.DB.Open(Session{}).Update(session)
}
