package service

import (
	"context"
	"time"

	"github.com/KnoblauchPilze/user-service/pkg/db"
	"github.com/KnoblauchPilze/user-service/pkg/persistence"
	"github.com/KnoblauchPilze/user-service/pkg/repositories"
	"github.com/google/uuid"
)

type mockUserRepository struct {
	repositories.UserRepository

	user      persistence.User
	ids       []uuid.UUID
	err       error
	updateErr error

	createCalled   int
	createdUser    persistence.User
	getCalled      int
	getId          uuid.UUID
	getEmailCalled int
	getEmail       string
	listCalled     int
	updateCalled   int
	updatedUser    persistence.User
	deleteCalled   int
	deleteId       uuid.UUID
}

type mockApiKeyRepository struct {
	repositories.ApiKeyRepository

	apiKey    persistence.ApiKey
	apiKeyIds []uuid.UUID
	createErr error
	getErr    error
	deleteErr error

	createCalled        int
	createdApiKey       persistence.ApiKey
	getForKeyCalled     int
	apiKeyId            uuid.UUID
	getForUserCalled    int
	getForUserTxCalled  int
	userId              uuid.UUID
	deleteCalled        int
	deleteIds           []uuid.UUID
	deleteForUserCalled int
	deleteUserId        uuid.UUID
}

type mockAclRepository struct {
	repositories.AclRepository

	aclIds []uuid.UUID
	acl    persistence.Acl

	getErr        error
	getForUserErr error
	deleteErr     error

	inAclIds         []uuid.UUID
	getCalled        int
	inUserId         uuid.UUID
	getForUserCalled int
	deleteCalled     int
}

type mockUserLimitRepository struct {
	repositories.UserLimitRepository

	userLimitIds []uuid.UUID
	userLimit    persistence.UserLimit

	getErr        error
	getForUserErr error
	deleteErr     error

	inUserLimitIds   []uuid.UUID
	getCalled        int
	inUserId         uuid.UUID
	getForUserCalled int
	deleteCalled     int
}

type mockUniverseRepository struct {
	repositories.UniverseRepository

	universe persistence.Universe
	err      error

	createCalled    int
	createdUniverse persistence.Universe
	getCalled       int
	getId           uuid.UUID
	listCalled      int
	deleteCalled    int
	deleteId        uuid.UUID
}

type mockPlanetRepository struct {
	repositories.PlanetRepository

	planet persistence.Planet
	err    error

	createCalled  int
	createdPlanet persistence.Planet
	getCalled     int
	getId         uuid.UUID
	listCalled    int
	deleteCalled  int
	deleteId      uuid.UUID
}

type mockPlayerRepository struct {
	repositories.PlayerRepository

	player persistence.Player
	err    error

	createCalled  int
	createdPlayer persistence.Player
	getCalled     int
	getId         uuid.UUID
	listCalled    int
	deleteCalled  int
	deleteId      uuid.UUID
}

type mockConnectionPool struct {
	db.ConnectionPool

	tx  mockTransaction
	err error
}

type mockTransaction struct {
	db.Transaction

	timeStamp time.Time

	closeCalled int
}

func (m *mockUserRepository) Create(ctx context.Context, user persistence.User) (persistence.User, error) {
	m.createCalled++
	m.createdUser = user
	return m.user, m.err
}

func (m *mockUserRepository) Get(ctx context.Context, id uuid.UUID) (persistence.User, error) {
	m.getCalled++
	m.getId = id
	return m.user, m.err
}

func (m *mockUserRepository) GetByEmail(ctx context.Context, email string) (persistence.User, error) {
	m.getEmailCalled++
	m.getEmail = email
	return m.user, m.err
}

func (m *mockUserRepository) List(ctx context.Context) ([]uuid.UUID, error) {
	m.listCalled++
	return m.ids, m.err
}

func (m *mockUserRepository) Update(ctx context.Context, user persistence.User) (persistence.User, error) {
	m.updateCalled++
	m.updatedUser = user
	return m.updatedUser, m.updateErr
}

func (m *mockUserRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

func (m *mockApiKeyRepository) Create(ctx context.Context, apiKey persistence.ApiKey) (persistence.ApiKey, error) {
	m.createCalled++
	m.createdApiKey = apiKey
	return apiKey, m.createErr
}

func (m *mockApiKeyRepository) GetForUser(ctx context.Context, user uuid.UUID) ([]uuid.UUID, error) {
	m.getForUserCalled++
	m.userId = user
	return m.apiKeyIds, m.getErr
}

func (m *mockApiKeyRepository) GetForUserTx(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	m.getForUserTxCalled++
	m.userId = user
	return m.apiKeyIds, m.getErr
}

func (m *mockApiKeyRepository) GetForKey(ctx context.Context, apiKey uuid.UUID) (persistence.ApiKey, error) {
	m.getForKeyCalled++
	m.apiKeyId = apiKey
	return m.apiKey, m.getErr
}

func (m *mockApiKeyRepository) Delete(ctx context.Context, ids []uuid.UUID) error {
	m.deleteCalled++
	m.deleteIds = ids
	return m.deleteErr
}

func (m *mockApiKeyRepository) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	m.deleteForUserCalled++
	m.deleteUserId = user
	return m.deleteErr
}

func (m *mockAclRepository) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.Acl, error) {
	m.getCalled++
	m.inAclIds = append(m.inAclIds, id)
	return m.acl, m.getErr
}

func (m *mockAclRepository) GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	m.getForUserCalled++
	m.inUserId = user
	return m.aclIds, m.getForUserErr
}

func (m *mockAclRepository) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	m.deleteCalled++
	m.inUserId = user
	return m.deleteErr
}

func (m *mockUserLimitRepository) Get(ctx context.Context, tx db.Transaction, id uuid.UUID) (persistence.UserLimit, error) {
	m.getCalled++
	m.inUserLimitIds = append(m.inUserLimitIds, id)
	return m.userLimit, m.getErr
}

func (m *mockUserLimitRepository) GetForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) ([]uuid.UUID, error) {
	m.getForUserCalled++
	m.inUserId = user
	return m.userLimitIds, m.getForUserErr
}

func (m *mockUserLimitRepository) DeleteForUser(ctx context.Context, tx db.Transaction, user uuid.UUID) error {
	m.deleteCalled++
	m.inUserId = user
	return m.deleteErr
}

func (m *mockUniverseRepository) Create(ctx context.Context, universe persistence.Universe) (persistence.Universe, error) {
	m.createCalled++
	m.createdUniverse = universe
	return m.universe, m.err
}

func (m *mockUniverseRepository) Get(ctx context.Context, id uuid.UUID) (persistence.Universe, error) {
	m.getCalled++
	m.getId = id
	return m.universe, m.err
}

func (m *mockUniverseRepository) List(ctx context.Context) ([]persistence.Universe, error) {
	m.listCalled++
	return []persistence.Universe{m.universe}, m.err
}

func (m *mockUniverseRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

func (m *mockPlanetRepository) Create(ctx context.Context, planet persistence.Planet) (persistence.Planet, error) {
	m.createCalled++
	m.createdPlanet = planet
	return m.planet, m.err
}

func (m *mockPlanetRepository) Get(ctx context.Context, id uuid.UUID) (persistence.Planet, error) {
	m.getCalled++
	m.getId = id
	return m.planet, m.err
}

func (m *mockPlanetRepository) List(ctx context.Context) ([]persistence.Planet, error) {
	m.listCalled++
	return []persistence.Planet{m.planet}, m.err
}

func (m *mockPlanetRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

func (m *mockPlayerRepository) Create(ctx context.Context, player persistence.Player) (persistence.Player, error) {
	m.createCalled++
	m.createdPlayer = player
	return m.player, m.err
}

func (m *mockPlayerRepository) Get(ctx context.Context, id uuid.UUID) (persistence.Player, error) {
	m.getCalled++
	m.getId = id
	return m.player, m.err
}

func (m *mockPlayerRepository) List(ctx context.Context) ([]persistence.Player, error) {
	m.listCalled++
	return []persistence.Player{m.player}, m.err
}

func (m *mockPlayerRepository) Delete(ctx context.Context, tx db.Transaction, id uuid.UUID) error {
	m.deleteCalled++
	m.deleteId = id
	return m.err
}

func (m *mockConnectionPool) StartTransaction(ctx context.Context) (db.Transaction, error) {
	return &m.tx, m.err
}

func (m *mockTransaction) Close(ctx context.Context) {
	m.closeCalled++
}

func (m *mockTransaction) TimeStamp() time.Time {
	return m.timeStamp
}