package users

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/aolus-software/clean-go-fiber/internal/config"
	apperrors "github.com/aolus-software/clean-go-fiber/internal/errors"
	"github.com/aolus-software/clean-go-fiber/internal/libs/cache"
	"github.com/aolus-software/clean-go-fiber/internal/libs/pagination"
	"github.com/aolus-software/clean-go-fiber/internal/libs/txmanager"
	"github.com/aolus-software/clean-go-fiber/internal/libs/utils"
	"github.com/aolus-software/clean-go-fiber/internal/models"
	"github.com/aolus-software/clean-go-fiber/internal/modules/auth"
	"github.com/aolus-software/clean-go-fiber/internal/repositories"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Service interface {
	List(ctx context.Context, q *pagination.Query) (*utils.PaginatedData[UserList], error)
	Create(ctx context.Context, actor *repositories.UserInformation, req *CreateUserRequest) error
	Detail(ctx context.Context, id uuid.UUID) (*UserDetail, error)
	Update(ctx context.Context, actor *repositories.UserInformation, id uuid.UUID, req *UpdateUserRequest) error
	ResetPassword(ctx context.Context, id uuid.UUID, req *ResetPasswordRequest) error
	SendVerificationEmail(ctx context.Context, id uuid.UUID) error
	SendResetPasswordEmail(ctx context.Context, id uuid.UUID) error
	Delete(ctx context.Context, actor *repositories.UserInformation, id uuid.UUID) error
}

type service struct {
	userRepo    repositories.UserRepository
	roleRepo    repositories.RoleRepository
	accountMail *auth.AccountMail
	tx          *txmanager.TxManager
	cache       *cache.Cache
	cfg         *config.Config
}

func NewService(
	userRepo repositories.UserRepository,
	roleRepo repositories.RoleRepository,
	accountMail *auth.AccountMail,
	tx *txmanager.TxManager,
	c *cache.Cache,
	cfg *config.Config,
) Service {
	return &service{userRepo: userRepo, roleRepo: roleRepo, accountMail: accountMail, tx: tx, cache: c, cfg: cfg}
}

func notFound() error {
	return apperrors.NewNotFoundError("User not found")
}

func emailTaken() error {
	return apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("email", "Email is already registered"))
}

func superuserOnly() error {
	return apperrors.NewForbiddenError("Only a superuser can manage superuser accounts")
}

func (s *service) List(ctx context.Context, q *pagination.Query) (*utils.PaginatedData[UserList], error) {
	users, total, err := s.userRepo.FindAll(ctx, q, s.cfg.App.Timezone)
	if err != nil {
		return nil, err
	}
	items := make([]UserList, 0, len(users))
	for _, u := range users {
		roles := make([]string, 0, len(u.Roles))
		for _, r := range u.Roles {
			roles = append(roles, r.Name)
		}
		items = append(items, UserList{
			ID:        u.ID.String(),
			Name:      u.Name,
			Email:     u.Email,
			Status:    u.Status,
			Remark:    u.Remark,
			Roles:     roles,
			CreatedAt: u.CreatedAt,
			UpdatedAt: u.UpdatedAt,
		})
	}
	return &utils.PaginatedData[UserList]{
		Data: items,
		Meta: utils.PaginationMeta{Page: q.Page, Limit: q.PerPage, TotalCount: total},
	}, nil
}

func (s *service) Create(ctx context.Context, actor *repositories.UserInformation, req *CreateUserRequest) error {
	email := strings.TrimSpace(req.Email)
	existing, err := s.userRepo.FindLiveByEmail(ctx, email, nil)
	if err != nil {
		return err
	}
	if existing != nil {
		return emailTaken()
	}
	roleIDs, err := s.resolveRoleIDs(ctx, actor, req.RoleIDs)
	if err != nil {
		return err
	}
	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return err
	}

	user := &models.User{Name: req.Name, Email: email, Password: hashed, Status: req.Status, Remark: req.Remark}
	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		repo := s.userRepo.WithTx(tx)
		if err := repo.Create(ctx, user); err != nil {
			return err
		}
		return repo.ReplaceRoles(ctx, user.ID, roleIDs)
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return emailTaken()
	}
	return err
}

func (s *service) Detail(ctx context.Context, id uuid.UUID) (*UserDetail, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, notFound()
	}
	roles := make([]RoleOption, 0, len(user.Roles))
	for _, r := range user.Roles {
		roles = append(roles, RoleOption{ID: r.ID.String(), Name: r.Name})
	}
	return &UserDetail{
		ID:              user.ID.String(),
		Name:            user.Name,
		Email:           user.Email,
		Status:          user.Status,
		Remark:          user.Remark,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Roles:           roles,
		CreatedAt:       user.CreatedAt,
		UpdatedAt:       user.UpdatedAt,
	}, nil
}

func (s *service) Update(ctx context.Context, actor *repositories.UserInformation, id uuid.UUID, req *UpdateUserRequest) error {
	user, err := s.findManageable(ctx, actor, id)
	if err != nil {
		return err
	}
	email := strings.TrimSpace(req.Email)
	existing, err := s.userRepo.FindLiveByEmail(ctx, email, &user.ID)
	if err != nil {
		return err
	}
	if existing != nil {
		return emailTaken()
	}
	roleIDs, err := s.resolveRoleIDs(ctx, actor, req.RoleIDs)
	if err != nil {
		return err
	}

	user.Name = req.Name
	user.Email = email
	user.Status = req.Status
	user.Remark = req.Remark
	err = s.tx.WithinTransaction(ctx, func(tx *gorm.DB) error {
		repo := s.userRepo.WithTx(tx)
		if err := repo.Update(ctx, user); err != nil {
			return err
		}
		return repo.ReplaceRoles(ctx, user.ID, roleIDs)
	})
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return emailTaken()
	}
	if err != nil {
		return err
	}
	if user.Status != models.UserStatusActive {
		return s.cache.RevokeSessions(ctx, user.ID)
	}
	return s.cache.ForgetUsers(ctx, user.ID)
}

func (s *service) ResetPassword(ctx context.Context, id uuid.UUID, req *ResetPasswordRequest) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return notFound()
	}
	hashed, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return err
	}
	user.Password = hashed
	if err := s.userRepo.Update(ctx, user); err != nil {
		return err
	}
	return s.cache.RevokeSessions(ctx, user.ID)
}

func (s *service) SendVerificationEmail(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return notFound()
	}
	if user.EmailVerifiedAt != nil {
		return apperrors.NewBadRequestError("Email already verified", apperrors.Field("email", "Email address is already verified"))
	}
	return s.accountMail.SendVerification(ctx, user)
}

func (s *service) SendResetPasswordEmail(ctx context.Context, id uuid.UUID) error {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return notFound()
	}
	return s.accountMail.SendPasswordReset(ctx, user)
}

func (s *service) Delete(ctx context.Context, actor *repositories.UserInformation, id uuid.UUID) error {
	user, err := s.findManageable(ctx, actor, id)
	if err != nil {
		return err
	}
	if err := s.userRepo.Delete(ctx, user.ID); err != nil {
		return err
	}
	return s.cache.RevokeSessions(ctx, user.ID)
}

func (s *service) findManageable(ctx context.Context, actor *repositories.UserInformation, id uuid.UUID) (*models.User, error) {
	user, err := s.userRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, notFound()
	}
	if !actor.HasRole(models.SuperuserRole) {
		for _, r := range user.Roles {
			if r.Name == models.SuperuserRole {
				return nil, superuserOnly()
			}
		}
	}
	return user, nil
}

func (s *service) resolveRoleIDs(ctx context.Context, actor *repositories.UserInformation, raw []string) ([]uuid.UUID, error) {
	ids := make([]uuid.UUID, 0, len(raw))
	for _, r := range raw {
		id, err := uuid.Parse(r)
		if err != nil {
			return nil, apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("roleIds", "must contain valid uuids"))
		}
		if !slices.Contains(ids, id) {
			ids = append(ids, id)
		}
	}
	count, err := s.roleRepo.CountByIDs(ctx, ids)
	if err != nil {
		return nil, err
	}
	if int(count) != len(ids) {
		return nil, apperrors.NewUnprocessableEntityError("Validation error", apperrors.Field("roleIds", "contains roles that do not exist"))
	}
	if !actor.HasRole(models.SuperuserRole) {
		superuser, err := s.roleRepo.FindByName(ctx, models.SuperuserRole, nil)
		if err != nil {
			return nil, err
		}
		if superuser != nil && slices.Contains(ids, superuser.ID) {
			return nil, superuserOnly()
		}
	}
	return ids, nil
}
